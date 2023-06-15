package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// workaround windows virus false positive

type WordData struct {
	frequency            int
	word                 string
	frequencyAcrossAll   float64
	frequencyAcrossValid float64
}

func (w WordData) String() string {
	return fmt.Sprintf("%d %v%% %v%% %s", w.frequency, w.frequencyAcrossAll, w.frequencyAcrossValid, w.word)
}

// small program to valid words and also calculate their (relative) frequency
// all errors are just panic'd upon for simplicity
func main() {
	lines := readLines()
	words := parseLines(lines)
	writeWords(words)
}

func readLines() []string {
	bytes, err := os.ReadFile("words.tsv")
	if err != nil {
		panic(err)
	}

	str := string(bytes)
	trimmed := strings.TrimSpace(str)
	return strings.Split(trimmed, "\n")
}

// gets the words that use letters only from the first half of the alphabet
func parseLines(lines []string) []WordData {
	words := make([]WordData, 0, int(math.Sqrt(float64(len(lines))))+1) // my estimate for the number of words that will be valid
	usagesAcrossAll := 0
	usagesAcrossValid := 0
	for _, line := range lines {
		parts := strings.Split(line, "\t")

		frequencyStr := parts[0]
		frequency, err := strconv.Atoi(frequencyStr)
		if err != nil {
			panic(err)
		}
		usagesAcrossAll += frequency

		word := parts[1]
		if isFirstHalfOfTheAlphabetOnly(word) {
			usagesAcrossValid += frequency

			wordStruct := WordData{
				frequency: frequency,
				word:      word,
			}

			words = append(words, wordStruct)
		}
	}

	for i := range words {
		words[i].frequencyAcrossAll = float64(words[i].frequency) / float64(usagesAcrossAll)
		words[i].frequencyAcrossValid = float64(words[i].frequency) / float64(usagesAcrossValid)
	}
	return words
}

func writeWords(words []WordData) {
	lines := make([]string, len(words), len(words))
	for i, word := range words {
		lines[i] = word.String()
	}

	err := os.WriteFile("out.txt", []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		panic(err)
	}
}

func isFirstHalfOfTheAlphabetOnly(word string) bool {
	for _, letter := range word {
		if !isFirstHalfOfTheAlphabet(letter) {
			return false
		}
	}
	return true
}

func isFirstHalfOfTheAlphabet(letter rune) bool {
	return letter >= 'a' && letter <= 'm'
}
