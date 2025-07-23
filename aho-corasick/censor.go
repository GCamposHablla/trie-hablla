package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/cloudflare/ahocorasick"
)

var (
		stringBuilderPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
)

type Censor struct {
	matcher *ahocorasick.Matcher
	words   []string
	wordSet map[string]bool 
}

func NewCensor(path string) (*Censor, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("erro abrindo lista: %w", err)
	}
	defer file.Close()

	var rawWords []string
	wordSet := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw := normalize(scanner.Text())
		if raw == "" {
			continue
		}
		if !wordSet[raw] {
			rawWords = append(rawWords, raw)
			wordSet[raw] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro lendo lista: %w", err)
	}

	keywords := make([][]byte, len(rawWords))
	for i, word := range rawWords {
		keywords[i] = []byte(word)
	}

	matcher := ahocorasick.NewMatcher(keywords)
	
	return &Censor{
		matcher: matcher,
		words:   rawWords,
		wordSet: wordSet,
	}, nil
}

func (c *Censor) Filter(input string) (found bool, censored string) {
	if input == "" {
		return false, input
	}

	normalized := normalize(input)
	hits := c.matcher.Match([]byte(normalized))
	
	if len(hits) == 0 {
		return false, input
	}

	return true, c.optimizedCensor(input, hits)
}

func (c *Censor) optimizedCensor(original string, hits []int) string {
	foundIndices := make(map[int]struct{}, len(hits))
	for _, idx := range hits {
		foundIndices[idx] = struct{}{}
	}

	result := stringBuilderPool.Get().(*strings.Builder)
	result.Reset()
	defer stringBuilderPool.Put(result)
	result.Grow(len(original))

	start := 0
	length := len(original)

	for start < length {
		for start < length && (original[start] == ' ' || original[start] == '\t' || original[start] == '\n') {
			result.WriteByte(original[start])
			start++
		}

		if start >= length {
			break
		}

		end := start
		for end < length && original[end] != ' ' && original[end] != '\t' && original[end] != '\n' {
			end++
		}

		word := original[start:end]

		normalizedWord := normalize(word)

		if c.wordSet[normalizedWord] {
			for i := 0; i < len(word); i++ {
				result.WriteByte('*')
			}
		} else {
			result.WriteString(word)
		}

		start = end
	}

	return result.String()
}