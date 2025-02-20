package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type wc struct {
	word  string
	count int
}

var validWord = regexp.MustCompile(`^[а-яА-Я-]+$`)

func Top10(t string) []string {
	const lastValue = 9
	s := strings.Fields(t)
	usedStack := make([]wc, 0)
	topWords := make([]string, 0)

	for _, v := range s {
		f := func(find string) (index int, err bool) {
			if !validWord.MatchString(find) {
				return 0, true
			}

			for i, mv := range usedStack {
				if find == mv.word {
					return i, false
				}
			}

			return 0, true
		}

		if v == "-" {
			continue
		}

		lw := strings.ReplaceAll(strings.ToLower(v), ".", "")
		lw = strings.ReplaceAll(lw, ",", "")

		index, err := f(lw)

		if err {
			usedStack = append(usedStack, wc{
				word:  lw,
				count: 1,
			})
			continue
		}

		usedStack[index].count++
	}

	sort.Slice(usedStack, func(i, j int) bool {
		if usedStack[i].count == usedStack[j].count {
			return usedStack[i].word < usedStack[j].word
		}

		return usedStack[i].count > usedStack[j].count
	})

	for index, v := range usedStack {
		topWords = append(topWords, v.word)

		print(v.word+": ", v.count, "\n")

		if index == lastValue {
			return topWords
		}
	}

	return topWords
}
