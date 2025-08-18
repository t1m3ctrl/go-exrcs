package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// Config хранит все параметры командной строки
type Config struct {
	after       int  // -A N
	before      int  // -B N
	context     int  // -C N
	count       bool // -c
	ignoreCase  bool // -i
	invert      bool // -v
	fixed       bool // -F
	lineNumbers bool // -n
	pattern     string
	filenames   []string
}

// Match представляет найденную строку вместе с контекстом
type Match struct {
	lineNum int
	content string
	isMatch bool // true — если совпадение, false — если контекст
}

// Grep описывает утилиту grep
type Grep struct {
	config  Config
	matcher func(string) bool
}

// NewGrep создаёт новый экземпляр grep
func NewGrep(config Config) (*Grep, error) {
	g := &Grep{config: config}

	// Настройка функции сопоставления
	pattern := config.pattern
	if config.ignoreCase {
		pattern = strings.ToLower(pattern)
	}

	if config.fixed {
		// Точное совпадение со строкой
		g.matcher = func(line string) bool {
			if config.ignoreCase {
				line = strings.ToLower(line)
			}
			return strings.Contains(line, pattern)
		}
	} else {
		// Совпадение с использованием регулярных выражений
		flags := ""
		if config.ignoreCase {
			flags = "(?i)"
		}
		re, err := regexp.Compile(flags + pattern)
		if err != nil {
			return nil, fmt.Errorf("некорректный шаблон регулярного выражения: %v", err)
		}
		g.matcher = re.MatchString
	}

	return g, nil
}

// processReader обрабатывает входные данные из reader
func (g *Grep) processReader(reader io.Reader, filename string) error {
	scanner := bufio.NewScanner(reader)
	var lines []string

	// Читаем все строки
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка чтения входных данных: %v", err)
	}

	// Ищем совпадения
	matches := g.findMatches(lines)

	// Выводим результат
	if g.config.count {
		fmt.Println(len(matches))
	} else {
		g.printMatches(matches, filename)
	}

	return nil
}

// findMatches находит все совпадения и строки контекста
func (g *Grep) findMatches(lines []string) []Match {
	var matches []Match
	contextBefore := g.config.before
	contextAfter := g.config.after

	// Обработка ключа -C (контекст вокруг совпадения)
	if g.config.context > 0 {
		contextBefore = g.config.context
		contextAfter = g.config.context
	}

	// Запоминаем уже добавленные строки, чтобы не было дубликатов
	processed := make(map[int]bool)

	for i, line := range lines {
		isMatch := g.matcher(line)

		// Инверсия совпадения при ключе -v
		if g.config.invert {
			isMatch = !isMatch
		}

		if isMatch {
			// Добавляем строки до совпадения
			start := i - contextBefore
			if start < 0 {
				start = 0
			}

			// Добавляем строки после совпадения
			end := i + contextAfter
			if end >= len(lines) {
				end = len(lines) - 1
			}

			// Добавляем все строки в диапазоне
			for j := start; j <= end; j++ {
				if !processed[j] {
					match := Match{
						lineNum: j + 1, // номера строк начинаются с 1
						content: lines[j],
						isMatch: j == i,
					}
					matches = append(matches, match)
					processed[j] = true
				}
			}
		}
	}

	// Если нужно только количество совпадений, фильтруем список
	if g.config.count {
		var actualMatches []Match
		for _, match := range matches {
			if match.isMatch {
				actualMatches = append(actualMatches, match)
			}
		}
		return actualMatches
	}

	return matches
}

// printMatches выводит найденные совпадения с форматированием
func (g *Grep) printMatches(matches []Match, filename string) {
	multipleFiles := len(g.config.filenames) > 1

	for i, match := range matches {
		var output strings.Builder

		// Если несколько файлов — добавляем имя файла
		if multipleFiles && filename != "" {
			output.WriteString(filename)
			output.WriteString(":")
		}

		// Добавляем номер строки, если указан ключ -n
		if g.config.lineNumbers {
			output.WriteString(fmt.Sprintf("%d:", match.lineNum))
		}

		// Добавляем содержимое строки
		output.WriteString(match.content)

		fmt.Println(output.String())

		// Добавляем разделитель между группами контекста (как в оригинальном grep)
		if i < len(matches)-1 {
			nextMatch := matches[i+1]
			// Если между строками есть разрыв — печатаем "--"
			if nextMatch.lineNum > match.lineNum+1 {
				fmt.Println("--")
			}
		}
	}
}

// processFile обрабатывает один файл
func (g *Grep) processFile(filename string) error {
	if filename == "-" || filename == "" {
		return g.processReader(os.Stdin, "")
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("невозможно открыть %s: %v", filename, err)
	}
	defer file.Close()

	return g.processReader(file, filename)
}

// Run запускает утилиту grep
func (g *Grep) Run() error {
	if len(g.config.filenames) == 0 {
		// Читаем из stdin
		return g.processFile("-")
	}

	// Обрабатываем каждый файл
	for _, filename := range g.config.filenames {
		if err := g.processFile(filename); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var config Config

	// Определяем флаги командной строки
	flag.IntVar(&config.after, "A", 0, "вывести N строк после совпадения")
	flag.IntVar(&config.before, "B", 0, "вывести N строк до совпадения")
	flag.IntVar(&config.context, "C", 0, "вывести N строк вокруг совпадения")
	flag.BoolVar(&config.count, "c", false, "вывести только количество совпадений")
	flag.BoolVar(&config.ignoreCase, "i", false, "игнорировать регистр")
	flag.BoolVar(&config.invert, "v", false, "инвертировать совпадение")
	flag.BoolVar(&config.fixed, "F", false, "искать фиксированную строку (без regexp)")
	flag.BoolVar(&config.lineNumbers, "n", false, "выводить номера строк")

	flag.Parse()

	// Получаем шаблон и список файлов
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Использование: %s [ОПЦИИ] ШАБЛОН [ФАЙЛ...]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	config.pattern = args[0]
	if len(args) > 1 {
		config.filenames = args[1:]
	}

	// Создаём и запускаем grep
	grep, err := NewGrep(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}

	if err := grep.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}
