package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// SortOptions содержит параметры сортировки
type SortOptions struct {
	column    int    // номер колонки для сортировки (0 = вся строка)
	numeric   bool   // числовая сортировка
	reverse   bool   // обратная сортировка
	unique    bool   // только уникальные строки
	separator string // разделитель колонок
}

// LineData представляет строку с её ключом для сортировки
type LineData struct {
	original string  // оригинальная строка
	sortKey  string  // ключ для сортировки
	numKey   float64 // числовой ключ (если используется -n)
}

// SortableLines реализует интерфейс sort.Interface
type SortableLines struct {
	lines   []LineData
	options SortOptions
}

func (s SortableLines) Len() int {
	return len(s.lines)
}

func (s SortableLines) Less(i, j int) bool {
	var result bool

	if s.options.numeric {
		// Числовая сортировка
		result = s.lines[i].numKey < s.lines[j].numKey
	} else {
		// Строковая сортировка
		result = s.lines[i].sortKey < s.lines[j].sortKey
	}

	// Обратная сортировка если нужно
	if s.options.reverse {
		result = !result
	}

	return result
}

func (s SortableLines) Swap(i, j int) {
	s.lines[i], s.lines[j] = s.lines[j], s.lines[i]
}

// extractSortKey извлекает ключ для сортировки из строки
func extractSortKey(line string, options SortOptions) (string, float64) {
	if options.column == 0 {
		// Сортировка по всей строке
		if options.numeric {
			if num, err := strconv.ParseFloat(strings.TrimSpace(line), 64); err == nil {
				return line, num
			}
			// Если не удалось парсить как число, используем 0
			return line, 0
		}
		return line, 0
	}

	// Сортировка по определённой колонке
	columns := strings.Split(line, options.separator)

	// Проверяем, что колонка существует (нумерация с 1)
	if options.column > len(columns) {
		// Если колонки нет, используем пустую строку
		if options.numeric {
			return "", 0
		}
		return "", 0
	}

	key := columns[options.column-1]

	if options.numeric {
		if num, err := strconv.ParseFloat(strings.TrimSpace(key), 64); err == nil {
			return key, num
		}
		// Если не удалось парсить как число, используем 0
		return key, 0
	}

	return key, 0
}

// readLines читает строки из reader
func readLines(reader io.Reader, options SortOptions) ([]LineData, error) {
	scanner := bufio.NewScanner(reader)
	var lines []LineData
	seenLines := make(map[string]bool) // для отслеживания уникальных строк

	for scanner.Scan() {
		line := scanner.Text()

		// Если нужны только уникальные строки, проверяем
		if options.unique {
			if seenLines[line] {
				continue // пропускаем дублирующуюся строку
			}
			seenLines[line] = true
		}

		sortKey, numKey := extractSortKey(line, options)

		lines = append(lines, LineData{
			original: line,
			sortKey:  sortKey,
			numKey:   numKey,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	return lines, nil
}

// processInput обрабатывает входные данные и выполняет сортировку
func processInput(reader io.Reader, writer io.Writer, options SortOptions) error {
	lines, err := readLines(reader, options)
	if err != nil {
		return err
	}

	// Создаём структуру для сортировки
	sortable := SortableLines{
		lines:   lines,
		options: options,
	}

	// Выполняем сортировку
	sort.Sort(sortable)

	// Выводим отсортированные строки
	for _, lineData := range sortable.lines {
		if _, err := fmt.Fprintln(writer, lineData.original); err != nil {
			return fmt.Errorf("error writing output: %w", err)
		}
	}

	return nil
}

func main() {
	// Определяем флаги командной строки
	var options SortOptions

	flag.IntVar(&options.column, "k", 0, "sort by column number N (1-indexed, 0 = whole line)")
	flag.BoolVar(&options.numeric, "n", false, "sort numerically")
	flag.BoolVar(&options.reverse, "r", false, "sort in reverse order")
	flag.BoolVar(&options.unique, "u", false, "output only unique lines")
	flag.StringVar(&options.separator, "t", "\t", "field separator (default: tab)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [FILE]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Sort lines from FILE or standard input.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s file.txt          # sort lines from file.txt\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -k 2 -n file.txt  # sort by 2nd column numerically\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -r -u             # reverse unique sort from stdin\n", os.Args[0])
	}

	flag.Parse()

	var reader io.Reader

	// Определяем источник данных (файл или stdin)
	if flag.NArg() > 0 {
		filename := flag.Arg(0)
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Error opening file %s: %v", filename, err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Printf("Error closing file: %v", err)
			}
		}()
		reader = file
	} else {
		reader = os.Stdin
	}

	// Обрабатываем входные данные
	if err := processInput(reader, os.Stdout, options); err != nil {
		log.Fatalf("Error processing input: %v", err)
	}
}
