package main

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	defaultVocabularyFilePath = "vocabulary.txt"
)

type MapStringToInt map[string]int
type MapIntToStringSlice map[int][]string

func (m *MapStringToInt) delete(word string) { delete(*m, word) }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func readInput(path string) (res MapStringToInt, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	res = make(map[string]int)
	for scanner.Scan() {
		word := strings.ToUpper(scanner.Text())
		if _, exist := res[word]; exist {
			res[word]++
		} else {
			res[word] = 1
		}
	}
	return res, scanner.Err()
}

func distanceLevenshtein(ctx *Context) int {
	var (
		s, t            = ctx.Task.value, ctx.word
		ls, lt          = len(s), len(t)
		column          = ctx.data[:ls+1]
		cost, last, old int
	)
	for y := 1; y <= ls; y++ {
		column[y] = y
	}
	for x := 1; x <= lt; x++ {
		column[0] = x
		last = x - 1
		for y := 1; y <= ls; y++ {
			old = column[y]
			cost = 0
			if s[y-1] != t[x-1] {
				cost = 1
			}
			column[y] = min(
				column[y]+1,
				min(column[y-1]+1,
					last+cost))
			last = old
		}
	}
	return column[ls]
}

type Task struct {
	value string
	count int
}

type Context struct {
	Task
	data []int
	word string
}

func newContext(size int) *Context { return &Context{data: make([]int, size)} }

type Vocabulary struct {
	MapIntToStringSlice
	maxWordLength int
}

func newVocabulary() *Vocabulary { return &Vocabulary{make(MapIntToStringSlice), 0} }

func (v *Vocabulary) loadFromFile(path string, process func(string)) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.ToUpper(scanner.Text())
		size := len(word)
		process(word)
		if size > v.maxWordLength {
			v.maxWordLength = size
		}
		v.MapIntToStringSlice[size] = append(v.MapIntToStringSlice[size], word)
	}
	return scanner.Err()
}

func (v *Vocabulary) min(ctx *Context) (res int) {
	res = len(ctx.Task.value)
	if num := v.searchInBucket(ctx, 0); num == 1 {
		return ctx.Task.count * 1
	} else {
		if num < res {
			res = num
		}
	}
	for offset := 1; offset < res; offset++ {
		if num := v.searchInBucket(ctx, offset); num == 1 {
			return ctx.Task.count * 1
		} else {
			if num < res {
				res = num
			}
		}
	}
	return ctx.Task.count * res
}

func (v *Vocabulary) searchInBucket(ctx *Context, offset int) int {
	var (
		len = len(ctx.Task.value)
		res = len
		num int
	)
	for _, ctx.word = range v.MapIntToStringSlice[len+offset] {
		num = distanceLevenshtein(ctx)
		if num == 1 {
			return 1
		}
		if num < res {
			res = num
		}
	}
	if offset == 0 {
		return res
	}
	for _, ctx.word = range v.MapIntToStringSlice[len-offset] {
		num = distanceLevenshtein(ctx)
		if num == 1 {
			return 1
		}
		if num < res {
			res = num
		}
	}
	return res
}

func getSumMin(words MapStringToInt, v *Vocabulary) (res uint64) {
	in, wg := make(chan Task, len(words)), new(sync.WaitGroup)
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			ctx := newContext(v.maxWordLength)
			for ctx.Task = range in {
				atomic.AddUint64(&res, uint64(v.min(ctx)))
			}
			wg.Done()
		}()
	}
	for word, count := range words {
		in <- Task{word, count}
	}
	close(in)
	wg.Wait()
	return
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Wrong number of arguments")
	}
	words, err := readInput(os.Args[1])
	if err != nil {
		log.Fatal("Read input error: ", err)
	}
	v := newVocabulary()
	err = v.loadFromFile(defaultVocabularyFilePath, words.delete)
	if err != nil {
		log.Fatal("Load vocabulary error: ", err)
	}
	println(getSumMin(words, v))
}
