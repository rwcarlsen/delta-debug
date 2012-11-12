package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"sort"
)

type Builder interface {
	BuildInput(Set) []byte
	Len() int
}

type Tester interface {
	Test(input []byte) Outcome
}

type word struct {
	words [][]byte
}

func ByWord(r io.Reader) (Builder, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	words := splitAny(data, "\n\t ")
	return &word{words: words}, nil
}

func (wb *word) BuildInput(set Set) []byte {
	sort.Ints([]int(set))

	inputwords := make([][]byte, len(set))
	for i, index := range set {
		inputwords[i] = wb.words[index]
	}

	input := bytes.Join(inputwords, []byte(""))
	return input
}

func (wb *word) Len() int {
	return len(wb.words)
}

type char struct {
	data []byte
}

func ByChar(r io.Reader) (Builder, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &char{data: data}, nil
}

func (cb *char) BuildInput(set Set) []byte {
	sort.Ints([]int(set))

	input := make([]byte, len(set))
	for i, index := range set {
		input[i] = cb.data[index]
	}

	return input
}

func (cb *char) Len() int {
	return len(cb.data)
}

type line struct {
	lines [][]byte
}

func ByLine(r io.Reader) (Builder, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(data, []byte("\n"))
	return &line{lines: lines}, nil
}

func (lb *line) BuildInput(set Set) []byte {
	sort.Ints([]int(set))

	inputlines := make([][]byte, len(set))
	for i, index := range set {
		inputlines[i] = lb.lines[index]
	}

	input := bytes.Join(inputlines, []byte("\n"))
	return append(input, byte('\n'))
}

func (lb *line) Len() int {
	return len(lb.lines)
}

type TestCase struct {
	T Tester
	B Builder
}

func (t *TestCase) Test(set Set) Outcome {
	input := t.B.BuildInput(set)
	return t.T.Test(input)
}

func (t *TestCase) Len() int {
	return t.B.Len()
}

// splitAny splits a byte slice at every occurence of any char in chars -
// collapsing consecutive char delimiters onto the preceding split chunk.
func splitAny(s []byte, chars string) [][]byte {
	beg := 0
	chunks := [][]byte{}
	for i, b := range s {
		for _, c := range chars {
			if b == byte(c) {
				end := i + 1
				nextNotWhite := bytes.IndexAny(s[end:end+1], chars) == -1
				if nextNotWhite {
					chunks = append(chunks, s[beg:end])
					beg = end
				}
				break
			}
		}
	}
	return append(chunks, s[beg:])
}
