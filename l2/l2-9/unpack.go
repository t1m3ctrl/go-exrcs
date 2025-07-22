package unpack

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// Unpack распаковывает строку с повторяющимися символами и поддержкой escape-последовательностей.
func Unpack(input string) (string, error) {
	var result strings.Builder
	runes := []rune(input)
	length := len(runes)

	var escape bool
	var prevRune rune
	var hasPrev bool

	for i := 0; i < length; i++ {
		curr := runes[i]

		switch {
		case escape:
			result.WriteRune(curr)
			prevRune = curr
			hasPrev = true
			escape = false

		case curr == '\\':
			escape = true

		case unicode.IsDigit(curr):
			if !hasPrev {
				return "", errors.New("invalid input: digit without preceding character")
			}
			j := i
			for j < length && unicode.IsDigit(runes[j]) {
				j++
			}
			countStr := string(runes[i:j])
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return "", err
			}
			if count == 0 {
				// удалить предыдущий символ
				resultStr := result.String()
				result.Reset()
				result.WriteString(resultStr[:len(resultStr)-len(string(prevRune))])
			} else {
				result.WriteString(strings.Repeat(string(prevRune), count-1))
			}
			i = j - 1 // пропускаем все цифры, уже обработаны
			hasPrev = false

		default:
			result.WriteRune(curr)
			prevRune = curr
			hasPrev = true
		}
	}

	if escape {
		return "", errors.New("invalid input: unfinished escape sequence")
	}

	return result.String(), nil
}
