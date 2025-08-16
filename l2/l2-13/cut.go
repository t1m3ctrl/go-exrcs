package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// parseFields парсит строку вида "1,3-5,7" в множество индексов (1-based)
func parseFields(spec string) map[int]struct{} {
	fields := make(map[int]struct{})
	parts := strings.Split(spec, ",")
	for _, p := range parts {
		if strings.Contains(p, "-") {
			bounds := strings.SplitN(p, "-", 2)
			start, err1 := strconv.Atoi(bounds[0])
			end, err2 := strconv.Atoi(bounds[1])
			if err1 == nil && err2 == nil && start > 0 && end >= start {
				for i := start; i <= end; i++ {
					fields[i] = struct{}{}
				}
			}
		} else {
			if num, err := strconv.Atoi(p); err == nil && num > 0 {
				fields[num] = struct{}{}
			}
		}
	}
	return fields
}

func main() {
	fieldSpec := flag.String("f", "", "fields (например: 1,3-5)")
	delimiter := flag.String("d", "\t", "delimiter (по умолчанию TAB)")
	separated := flag.Bool("s", false, "show only lines with delimiter")
	flag.Parse()

	if *fieldSpec == "" {
		fmt.Fprintln(os.Stderr, "Ошибка: нужно указать -f")
		os.Exit(1)
	}

	fields := parseFields(*fieldSpec)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		if *separated && !strings.Contains(line, *delimiter) {
			continue
		}

		cols := strings.Split(line, *delimiter)
		var out []string

		for i := 1; i <= len(cols); i++ {
			if _, ok := fields[i]; ok {
				out = append(out, cols[i-1])
			}
		}

		if len(out) > 0 {
			fmt.Println(strings.Join(out, *delimiter))
		} else {
			fmt.Println()
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения:", err)
		os.Exit(1)
	}
}
