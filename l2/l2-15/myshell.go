package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Сигнатура функции встроенной команды
// stdin/stdout позволяют встроенным командам участвовать в конвейерах
type builtinFunc func(args []string, stdin io.Reader, stdout io.Writer) error

// activeProcs отслеживает текущие выполняющиеся внешние команды, чтобы мы могли прервать их при Ctrl+C.
var activeProcs = struct {
	cmds []*exec.Cmd
}{cmds: []*exec.Cmd{}}

func main() {
	fmt.Println("my-shell started")

	// Обработка Ctrl+C (SIGINT)
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)
	go func() {
		for range interruptCh {
			// Убить все активные процессы в текущем конвейере
			for _, c := range activeProcs.cmds {
				if c != nil && c.Process != nil {
					_ = killProcess(c)
				}
			}
			// Очистить список после прерывания
			activeProcs.cmds = nil
			fmt.Println() // перейти на новую строку после ^C
		}
	}()

	r := bufio.NewReader(os.Stdin)
	for {
		cwd, _ := os.Getwd()
		fmt.Printf("%s> ", filepath.Base(cwd))
		line, err := r.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("\nbye!")
				return
			}
			fmt.Fprintf(os.Stderr, "ошибка чтения: %v\n", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if err := runLine(line); err != nil {
			fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
		}
	}
}

func runLine(line string) error {
	segments := splitPipeline(line)
	if len(segments) == 0 {
		return nil
	}

	// Особый случай: одиночные встроенные команды, которые изменяют состояние оболочки, не должны использоваться в пайплайне
	if len(segments) == 1 {
		cmd, args := parseCmd(segments[0])
		if cmd == "cd" {
			return builtinCd(args)
		}
		if cmd == "kill" {
			return builtinKill(args)
		}
	}

	// Построение пайплайна из команд
	var cmds []*exec.Cmd         // внешние команды
	var pipes []io.ReadCloser    // каналы для чтения
	var writers []io.WriteCloser // каналы для записи

	// Соединяем этапы один за другим.
	var prevReader io.Reader
	for i, seg := range segments {
		cmdName, args := parseCmd(seg)

		// Создаем пайп для передачи вывода между этапами, если это не последний этап
		var stdoutW io.WriteCloser
		if i < len(segments)-1 {
			pr, pw := io.Pipe()
			pipes = append(pipes, pr)
			writers = append(writers, pw)
			stdoutW = pw
		} else {
			stdoutW = nopWriteCloser{os.Stdout}
		}

		stdinR := prevReader
		if stdinR == nil {
			stdinR = os.Stdin
		}

		// Обработка встроенных команд, которые могут быть частью пайплайна
		if isBuiltinProcessLike(cmdName) {
			// Запуск встроенного этапа сразу и передача вывода следующему через горутину (для потокового поведения)
			pr, pw := io.Pipe()
			// Если есть upstream reader, игнорируем его для echo/pwd/ps; для будущей совместимости передаем.
			go func(name string, a []string, in io.Reader, out io.Writer) {
				defer func() { _ = pw.Close() }()
				// Записать вывод встроенной команды в pw; перенаправить в stdoutW, если это последний этап
				var target io.Writer = pw
				if i == len(segments)-1 {
					target = stdoutW
				}
				_ = runBuiltinProcessLike(name, a, in, target)
			}(cmdName, args, stdinR, stdoutW)

			prevReader = pr
			continue
		}

		// Внешняя команда
		c := exec.Command(cmdName, args...)

		c.Stdin = stdinR
		if i < len(segments)-1 {
			c.Stdout = stdoutW
		} else {
			c.Stdout = os.Stdout
		}
		c.Stderr = os.Stderr

		cmds = append(cmds, c)
		prevReader = pipesReaderFor(i, pipes)
	}

	// Соединяем каналы между внешними командами
	for i, c := range cmds {
		if i > 0 {
			c.Stdin = pipes[i-1]
		}
	}

	// Запуск всех внешних команд
	activeProcs.cmds = nil
	for _, c := range cmds {
		if err := c.Start(); err != nil {
			// Закрываем writer-концы, чтобы избежать зависаний
			for _, w := range writers {
				_ = w.Close()
			}
			return fmt.Errorf("не удалось запустить %s: %w", c.Path, err)
		}
		activeProcs.cmds = append(activeProcs.cmds, c)
	}

	// Закрываем writer-концы в родителе, чтобы downstream читатели увидели EOF
	for _, w := range writers {
		_ = w.Close()
	}

	// Ждем завершения всех команд
	var waitErr error
	for _, c := range cmds {
		if err := c.Wait(); err != nil {
			waitErr = err
		}
	}
	activeProcs.cmds = nil

	return waitErr
}

func splitPipeline(line string) []string {
	parts := strings.Split(line, "|")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func parseCmd(seg string) (string, []string) {
	fields := strings.Fields(seg)
	if len(fields) == 0 {
		return "", nil
	}
	return fields[0], fields[1:]
}

func isBuiltinProcessLike(name string) bool {
	switch name {
	case "echo", "pwd", "ps", "ls":
		return true
	default:
		return false
	}
}

func runBuiltinProcessLike(name string, args []string, stdin io.Reader, stdout io.Writer) error {
	switch name {
	case "echo":
		_, err := fmt.Fprintln(stdout, strings.Join(args, " "))
		return err
	case "pwd":
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(stdout, wd)
		return err
	case "ps":
		return builtinPs(stdout)
	case "ls":
		return builtinLs(args, stdout)
	default:
		return fmt.Errorf("неизвестная встроенная команда: %s", name)
	}
}

// cd изменяет рабочую директорию оболочки (не для пайплайнов)
func builtinCd(args []string) error {
	if len(args) == 0 {
		// по умолчанию HOME / USERPROFILE
		home := os.Getenv("HOME")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		if home == "" {
			return errors.New("cd: HOME не задан")
		}
		return os.Chdir(home)
	}
	path := args[0]
	return os.Chdir(path)
}

// kill отправляет сигнал завершения процессу (кроссплатформенно ?? не тестировалось на UNIX подобных системах)
func builtinKill(args []string) error {
	if len(args) < 1 {
		return errors.New("использование: kill <pid>")
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("kill: неверный pid: %v", err)
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("kill: %v", err)
	}
	if runtime.GOOS == "windows" {
		// используем Kill (TerminateProcess)
		return proc.Kill()
	}
	// На Unix сначала пытаемся SIGTERM
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		// Фоллбэк на SIGKILL
		_ = proc.Signal(syscall.SIGKILL)
		return err
	}
	return nil
}

// ps выводит список процессов
func builtinPs(w io.Writer) error {
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("powershell", "-Command", "Get-Process | Format-Table -AutoSize")
	} else {
		c = exec.Command("ps", "aux")
	}
	c.Stdout = w
	c.Stderr = os.Stderr
	return c.Run()
}

func builtinLs(args []string, w io.Writer) error {
	var dir string
	if len(args) > 0 {
		dir = args[0]
	} else {
		dir = "."
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	fmt.Println()
	for _, e := range entries {
		if e.IsDir() {
			// Синий цвет для директорий (ANSI escape код)
			fmt.Fprintln(w, "\033[34m", e.Name(), "\033[0m")
		} else {
			fmt.Fprintln(w, e.Name())
		}
	}
	return nil
}

// вспомогательные функции

type nopWriteCloser struct{ io.Writer }

func (n nopWriteCloser) Close() error { return nil }

func pipesReaderFor(i int, pipes []io.ReadCloser) io.Reader {
	if i == 0 {
		return nil
	}
	return pipes[i-1]
}

func killProcess(c *exec.Cmd) error {
	if c == nil || c.Process == nil {
		return nil
	}
	// Попытка корректного завершения на Unix, затем принудительно; на Windows Kill.
	if runtime.GOOS == "windows" {
		return c.Process.Kill()
	}
	// Сначала отправляем SIGINT для имитации Ctrl+C
	_ = c.Process.Signal(os.Interrupt)
	time.Sleep(100 * time.Millisecond)
	// Затем SIGTERM
	_ = c.Process.Signal(syscall.SIGTERM)
	time.Sleep(200 * time.Millisecond)
	// И наконец SIGKILL, если все еще выполняется
	return c.Process.Kill()
}
