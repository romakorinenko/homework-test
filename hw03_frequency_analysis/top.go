package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Word struct {
	name  string
	count int
}

func Top10(input string) []string {
	splittedInput := strings.Fields(input)
	wordToCount := getWordsToCount(splittedInput)
	words := getWordsSlice(wordToCount)
	sortWordsSlice(words)

	return getResultWordSlice(words)
}

func getResultWordSlice(words []Word) []string {
	result := make([]string, 0, 10)
	for i := 0; i < len(words); i++ {
		if len(result) == 10 {
			return result
		}
		result = append(result, words[i].name)
	}
	return result
}

func sortWordsSlice(words []Word) {
	sort.Slice(words, func(i, j int) bool {
		if words[i].count == words[j].count {
			return words[i].name < words[j].name
		}
		return words[i].count > words[j].count
	})
}

func getWordsSlice(wordToCount map[string]int) []Word {
	words := make([]Word, 0)
	for name, count := range wordToCount {
		word := Word{name: name, count: count}
		words = append(words, word)
	}
	return words
}

func getWordsToCount(input []string) map[string]int {
	wordToCount := make(map[string]int)
	for _, word := range input {
		if word == "" {
			continue
		}
		wordToCount[word]++
	}
	return wordToCount
}
