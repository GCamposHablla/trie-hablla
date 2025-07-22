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
	// Pool de buffers para reduzir alocações
	stringBuilderPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
)

type Censor struct {
	matcher *ahocorasick.Matcher
	words   []string
	wordSet map[string]bool // Para lookup O(1)
}

// Carrega palavras do arquivo e constrói automato
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
		if !wordSet[raw] { // Evita duplicatas
			rawWords = append(rawWords, raw)
			wordSet[raw] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro lendo lista: %w", err)
	}

	// Converte para [][]byte uma única vez
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

// Verifica e censura palavras encontradas - VERSÃO OTIMIZADA
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

// Nova função de censura otimizada
func (c *Censor) optimizedCensor(original string, hits []int) string {
	// Cria mapa de palavras encontradas para lookup O(1)
	foundWords := make(map[string]bool, len(hits))
	for _, idx := range hits {
		foundWords[c.words[idx]] = true
	}

	result := stringBuilderPool.Get().(*strings.Builder)
	result.Reset()
	defer stringBuilderPool.Put(result)
	result.Grow(len(original))

	words := strings.Fields(original)
	for i, word := range words {
		if i > 0 {
			result.WriteByte(' ')
		}

		normalizedWord := normalize(word)
		if foundWords[normalizedWord] {
			result.WriteString(strings.Repeat("*", len(word)))
		} else {
			result.WriteString(word)
		}
	}

	return result.String()
}