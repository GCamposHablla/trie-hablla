package main

import (
	"strings"
	"sync"
	"unicode"
	"golang.org/x/text/unicode/norm"
)

var (
	leetMap = map[rune]rune{
		'4': 'a', '@': 'a', '^': 'a',
		'8': 'b',
		'<': 'c', '(': 'c',
		')': 'd',
		'3': 'e', '&': 'e', '€': 'e',
		'6': 'g',
		'1': 'i', '!': 'i',
		'0': 'o',
		'5': 's', '$': 's',
		'7': 't', '+': 't',
		'2': 'z',
	}

	// Pool para reutilizar strings.Builder
	normBuilderPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
)

// Versão otimizada da normalização
func normalize(input string) string {
	if input == "" {
		return ""
	}

	builder := normBuilderPool.Get().(*strings.Builder)
	builder.Reset()
	defer normBuilderPool.Put(builder)

	// Pre-aloca capacidade estimada
	builder.Grow(len(input))

	// Normaliza NFD uma única vez
	normalized := norm.NFD.String(input)
	
	for _, r := range normalized {
		// Skip diacríticos (acentos)
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		
		// Mapeia leet speak
		if mapped, ok := leetMap[r]; ok {
			r = mapped
		}
		
		// Apenas letras e números, em lowercase
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			builder.WriteRune(unicode.ToLower(r))
		}
	}
	
	return builder.String()
}