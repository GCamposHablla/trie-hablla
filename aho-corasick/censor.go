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

	// Cache de normalização para evitar reprocessamento
	normalizationCache = sync.Map{}

	// Pool para slices de palavras reutilizáveis
	wordsSlicePool = sync.Pool{
		New: func() interface{} {
			return make([]string, 0, 50) // capacidade inicial razoável
		},
	}
)

type Censor struct {
	matcher *ahocorasick.Matcher
	words   []string
	wordSet map[string]bool
}

type MatchRange struct {
	Start, End int
	Word       string
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

// normalizeWithCache usa cache para evitar renormalizar strings repetidas
func (c *Censor) normalizeWithCache(input string) string {
	if cached, ok := normalizationCache.Load(input); ok {
		return cached.(string)
	}

	result := normalize(input)

	// Limitar cache para evitar vazamento de memória
	// Política simples: só cachear strings pequenas e médias
	if len(input) <= 100 {
		normalizationCache.Store(input, result)
	}

	return result
}

func (c *Censor) Filter(input string) (found bool, censored string) {
	if input == "" {
		return false, input
	}

	normalized := c.normalizeWithCache(input)
	hits := c.matcher.Match([]byte(normalized))

	if len(hits) == 0 {
		return false, input
	}

	return true, c.applyCensoringByMatches(input, normalized, hits)
}

// buildFoundWords constrói set de palavras encontradas pelos matches
func (c *Censor) buildFoundWords(hits []int) map[string]bool {
	foundWords := make(map[string]bool, len(hits))
	for _, idx := range hits {
		if idx >= 0 && idx < len(c.words) {
			foundWords[c.words[idx]] = true
		}
	}
	return foundWords
}

// splitWordsOptimized divide o texto em palavras de forma otimizada
func (c *Censor) splitWordsOptimized(text string) []string {
	words := wordsSlicePool.Get().([]string)
	words = words[:0] // reset mas mantém capacidade
	defer wordsSlicePool.Put(words)

	start := -1
	result := make([]string, 0, len(words))

	for i, r := range text {
		isSpace := r == ' ' || r == '\t' || r == '\n' || r == '\r'

		if !isSpace && start == -1 {
			start = i
		} else if isSpace && start != -1 {
			result = append(result, text[start:i])
			start = -1
		}
	}

	// Capturar última palavra se não termina com espaço
	if start != -1 {
		result = append(result, text[start:])
	}

	return result
}

// applyCensoringByMatches aplica censura baseada nos matches do Aho-Corasick
func (c *Censor) applyCensoringByMatches(original, normalized string, hits []int) string {
	if len(hits) == 0 {
		return original
	}

	// Construir set de palavras encontradas
	foundWords := c.buildFoundWords(hits)

	result := stringBuilderPool.Get().(*strings.Builder)
	result.Reset()
	defer stringBuilderPool.Put(result)
	result.Grow(len(original) + len(original)/4) // crescer um pouco mais para acomodar possíveis expansões

	words := c.splitWordsOptimized(original)

	// Usar índices para reconstruir com espaçamento original
	wordIndex := 0
	i := 0

	for i < len(original) {
		r := rune(original[i])

		// Escrever espaços/tabs/newlines diretamente
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			result.WriteByte(original[i])
			i++
			continue
		}

		// Encontrar fim da palavra atual
		start := i
		for i < len(original) {
			r := rune(original[i])
			if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
				break
			}
			i++
		}

		word := original[start:i]

		if wordIndex < len(words) {
			// Usar cache de normalização apenas para a palavra
			normalizedWord := c.normalizeWithCache(word)

			if foundWords[normalizedWord] {
				// Otimização: usar strings.Repeat em vez de loop
				result.WriteString(strings.Repeat("*", len(word)))
			} else {
				result.WriteString(word)
			}
			wordIndex++
		} else {
			// Fallback: escrever palavra original
			result.WriteString(word)
		}
	}

	return result.String()
}

// Método adicional para estatísticas (útil para debugging/monitoramento)
func (c *Censor) GetCacheStats() (int, int) {
	hits := 0
	misses := 0

	// Simples contagem aproximada do cache
	normalizationCache.Range(func(key, value interface{}) bool {
		hits++
		return true
	})

	return hits, misses
}

// ClearCache limpa o cache de normalização (útil para testes ou gestão de memória)
func (c *Censor) ClearCache() {
	normalizationCache.Range(func(key, value interface{}) bool {
		normalizationCache.Delete(key)
		return true
	})
}
