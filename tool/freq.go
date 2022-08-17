package main

import (
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/huichen/sego"
)

type WordFreq struct {
	Word     string
	Count    int
	Freq     float64
	Disperse float64
	Rank     int
	Score    float64
}

const (
	textPath = "../text/"
)

var lemma, _ = ioutil.ReadFile(textPath + "lemmatization.txt")
var segmenter sego.Segmenter

func main() {
	batch("cn", "freqChar")
	segmenter.LoadDictionary(textPath + "segment.txt")
	batch("cn", "freqWord")
	batch("en", "freqEnglish")
}

func loadScope() map[string][]byte {
	scope := map[string][]byte{}
	scope["freqChar"], _ = ioutil.ReadFile(textPath + "8105.txt")
	scope["freqWord"], _ = ioutil.ReadFile(textPath + "13436.txt")
	scope["freqEnglish"], _ = ioutil.ReadFile(textPath + "20000.txt")
	return scope
}

func batch(bookPath, proc string) {
	scope := loadScope()
	bookPath = textPath + bookPath + "/"
	files, _ := ioutil.ReadDir(bookPath)

	call := map[string]func(a, b []byte) []WordFreq{"freqChar": freqChar, "freqWord": freqWord, "freqEnglish": freqEnglish}

	scoreMap := map[string]float64{}
	for i, file := range files {
		println(proc, i, file.Name())
		book, _ := ioutil.ReadFile(bookPath + file.Name())

		wf := call[proc](book, scope[proc])

		for _, v := range wf {
			scoreMap[v.Word] += v.Score
		}
	}

	score := []WordFreq{}
	for k, v := range scoreMap {
		score = append(score, WordFreq{k, 0, 0, 0, 0, v})
	}
	score = sortWord(score, 5, "desc")

	output(score, textPath+proc+".txt")
}

func output(wf []WordFreq, fileName string) {
	res := ""
	for _, v := range wf {
		res += fmt.Sprintf("%s\t%d\t%f\t%f\t%d\t%f\n", v.Word, v.Count, v.Freq, v.Disperse, v.Rank, v.Score)
	}
	ioutil.WriteFile(fileName, []byte(res), 0644)
}

func sortWord(wf []WordFreq, col int, order string) []WordFreq {
	sort.Slice(wf, func(i, j int) bool {
		if col == 1 {
			if order == "desc" {
				return wf[i].Count > wf[j].Count
			} else {
				return wf[i].Count < wf[j].Count
			}
		} else if col == 2 {
			if order == "desc" {
				return wf[i].Freq > wf[j].Freq
			} else {
				return wf[i].Freq < wf[j].Freq
			}
		} else if col == 3 {
			if order == "desc" {
				return wf[i].Disperse > wf[j].Disperse
			} else {
				return wf[i].Disperse < wf[j].Disperse
			}
		} else if col == 4 {
			if order == "desc" {
				return wf[i].Rank > wf[j].Rank
			} else {
				return wf[i].Rank < wf[j].Rank
			}
		} else {
			if order == "desc" {
				return wf[i].Score > wf[j].Score
			} else {
				return wf[i].Score < wf[j].Score
			}
		}
	})
	return wf
}

func count[S rune | string](text []S) map[S]int {
	wc := map[S]int{}
	for _, v := range text {
		wc[v]++
	}
	return wc
}

func filter[N int | float64, S rune | string](wc map[S]N, std []S) map[S]N {
	intersect := map[S]int{}
	for _, v := range std {
		intersect[v]++
	}
	for k := range wc {
		intersect[k]++
	}
	for k, v := range intersect {
		if v < 2 {
			delete(wc, k)
		}
	}
	return wc
}
