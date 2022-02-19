package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Freq struct {
	Word      string
	Frequency int
}

func (f *Freq) Match(pattern, letters, badLetters string) bool {
	matched, err := regexp.MatchString(pattern, f.Word)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if !matched {
		return false
	}
	return containsLetters(f.Word, letters, badLetters)
}
func containsLetters(word, letters, badLetters string) bool {
	for _, r := range badLetters {
		if strings.ContainsRune(word, r) {
			return false
		}
	}

	wordRunes := []rune(word)
	for _, knownLetter := range letters {
		hasLetter := false
		for s, wordLetter := range wordRunes {
			if knownLetter != wordLetter {
				continue
			}
			hasLetter = true

			wordRunes = append(wordRunes[:s], wordRunes[s+1:]...)
			break
		}
		if !hasLetter {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println(usage())
		return
	}
	wordPattern, extraLetters, badLetters := parseArgs()

	freqs, err := os.Open("unigram_freq.csv")
	if err != nil {
		log.Fatalf("Can't open unigram file %s", err.Error())
	}
	defer freqs.Close()
	r := csv.NewReader(freqs)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Can't parse csv %s", err.Error())
	}

	wordFreqs := make([]Freq, 0)
	for _, rec := range records {
		if len(rec[0]) != 5 {
			continue
		}
		wordCount, err := strconv.Atoi(rec[1])
		if err != nil {
			continue
		}
		nf := Freq{
			Word:      rec[0],
			Frequency: wordCount,
		}
		wordFreqs = append(wordFreqs, nf)
	}

	for _, f := range wordFreqs {
		if f.Match(wordPattern, extraLetters, badLetters) {
			fmt.Println(f.Word)
		}
	}
}

func usage() string {
	return "Call wordle with two arguments. One to indicate known placements and the other to indicate known letters : ./wordle __e__ y"
}

func parseArgs() (string, string, string) {
	slate := os.Args[1]
	if len(slate) != 5 {
		log.Fatalf("Please input five letters or underscores to indicate the current slate")
	}

	letters := os.Args[2]

	knownLetters := strings.Replace(slate, "_", "", 5)
	knownLetters = fmt.Sprintf("%s%s", knownLetters, letters)
	if len(knownLetters) > 5 {
		log.Fatalf("More than give known letters.")
	}

	badLetters := os.Args[3]

	reLetter := "[a-z]"
	reBody := ""

	for _, letter := range slate {
		if letter == '_' {
			reBody = reBody + reLetter
		} else {
			reBody = reBody + string(letter)
		}
	}

	return reBody, letters, badLetters
}
