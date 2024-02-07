package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(raw_text string) []string {
	words_list := count(raw_text)
	rate_list := rateBy(words_list)

	return ret_10_words(rate_list, 10)
}

func count(text string) map[string]int {
	words := strings.Fields(text)
	wordCount := make(map[string]int)

	for _, word := range words {
		_, exists := wordCount[word]
		if exists {
			wordCount[word] += 1
		} else {
			wordCount[word] = 1
		}
	}
	return wordCount
}

func rateBy(words map[string]int) pairs {
	pairs := make(pairs, len(words))
	i := 0
	for k, v := range words {
		pairs[i] = pair{k, v}
		i++
	}

	sort.Sort(pairs)
	return pairs
}

func ret_10_words(p pairs, n int) []string {
	if len(p) == 0 {
		return []string{}
	}
	top := make([]string, 0)
	for i := 0; i < len(p) && i < n; i++ {
		top = append(top, p[i].word)
	}
	return top
}

type pair struct {
	word  string
	count int
}

type pairs []pair

func (p pairs) Len() int {
	return len(p)
}

func (p pairs) Less(i, j int) bool {
	if p[i].count == p[j].count {
		return p[i].word < p[j].word
	}
	return p[i].count > p[j].count
}

func (p pairs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

