package main

import (
	"sort"
	"strings"
)

func sortString(w string) string {
	runes := []rune(w)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}

func anagram(words []string) map[string][]string {
	// groups хранит для каждой отсортированной строки множество соответствующих слов
	groups := make(map[string]map[string]bool)
	// firstWord запоминает первое слово для каждой группы анаграмм
	firstWord := make(map[string]string)

	for _, word := range words {
		wordLower := strings.ToLower(word)
		sorted := sortString(wordLower)

		// Создаем новую группу, если такой отсортированной строки еще не было
		if _, exists := groups[sorted]; !exists {
			groups[sorted] = make(map[string]bool)
			firstWord[sorted] = wordLower
		}
		groups[sorted][wordLower] = true
	}

	// Формирование результата
	result := make(map[string][]string)
	for sortedKey, wordSet := range groups {
		if len(wordSet) <= 1 {
			continue
		}
		rep := firstWord[sortedKey]
		wordsSlice := make([]string, 0, len(wordSet))
		for word := range wordSet {
			wordsSlice = append(wordsSlice, word)
		}
		sort.Strings(wordsSlice)
		result[rep] = wordsSlice
	}
	return result
}

func main() {
	words := []string{"пяТак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	anagrams := anagram(words)
	for key, list := range anagrams {
		println(key, ":", strings.Join(list, ", "))
	}
}
