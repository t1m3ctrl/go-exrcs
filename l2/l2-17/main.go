package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type TelnetClient struct {
	host    string
	port    string
	timeout time.Duration
	conn    net.Conn
}

func NewTelnetClient(host, port string, timeout time.Duration) *TelnetClient {
	return &TelnetClient{
		host:    host,
		port:    port,
		timeout: timeout,
	}
}

func (tc *TelnetClient) Connect() error {
	address := net.JoinHostPort(tc.host, tc.port)

	// Создаем контекст с таймаутом для подключения
	ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
	defer cancel()

	// Используем DialContext для подключения с таймаутом
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к %s: %w", address, err)
	}

	tc.conn = conn
	fmt.Printf("Подключено к %s\n", address)
	return nil
}

func (tc *TelnetClient) Close() {
	if tc.conn != nil {
		tc.conn.Close()
		fmt.Println("Соединение закрыто")
	}
}

// readFromServer читает данные из сокета и выводит в STDOUT
func (tc *TelnetClient) readFromServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	reader := bufio.NewReader(tc.conn)
	buffer := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Устанавливаем небольшой таймаут для чтения, чтобы периодически проверять контекст
			tc.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

			n, err := reader.Read(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// Таймаут чтения - продолжаем цикл
					continue
				}
				if err == io.EOF {
					fmt.Println("\nСервер закрыл соединение")
				} else {
					fmt.Printf("\nОшибка чтения из сокета: %v\n", err)
				}
				return
			}

			if n > 0 {
				os.Stdout.Write(buffer[:n])
			}
		}
	}
}

// writeToServer читает из STDIN и отправляет в сокет
func (tc *TelnetClient) writeToServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Создаем канал для получения результата сканирования
			inputChan := make(chan bool, 1)

			go func() {
				inputChan <- scanner.Scan()
			}()

			// Ждем либо ввод, либо отмену контекста
			select {
			case <-ctx.Done():
				return
			case hasInput := <-inputChan:
				if !hasInput {
					// EOF (Ctrl+D) или ошибка
					if err := scanner.Err(); err != nil {
						fmt.Printf("Ошибка чтения из STDIN: %v\n", err)
					} else {
						fmt.Println("\nПолучен EOF (Ctrl+D), завершение работы...")
					}
					return
				}

				text := scanner.Text() + "\n"
				_, err := tc.conn.Write([]byte(text))
				if err != nil {
					fmt.Printf("Ошибка записи в сокет: %v\n", err)
					return
				}
			}
		}
	}
}

func (tc *TelnetClient) Run() error {
	// Создаем контекст для управления горутинами
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обрабатываем сигналы для корректного завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Запускаем горутину для чтения из сервера
	wg.Add(1)
	go tc.readFromServer(ctx, &wg)

	// Запускаем горутину для записи на сервер
	wg.Add(1)
	go tc.writeToServer(ctx, &wg)

	// Горутина для обработки сигналов
	go func() {
		select {
		case sig := <-sigChan:
			fmt.Printf("\nПолучен сигнал %v, завершение работы...\n", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	// Ждем завершения всех горутин
	wg.Wait()

	return nil
}

func main() {
	var (
		timeout = flag.Duration("timeout", 10*time.Second, "Таймаут подключения")
	)
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Использование: %s [--timeout=10s] <хост> <порт>\n", os.Args[0])
		os.Exit(1)
	}

	host := args[0]
	port := args[1]

	client := NewTelnetClient(host, port, *timeout)
	defer client.Close()

	// Подключаемся к серверу
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка подключения: %v\n", err)
		os.Exit(1)
	}

	// Запускаем основной цикл обработки
	if err := client.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка выполнения: %v\n", err)
		os.Exit(1)
	}
}
