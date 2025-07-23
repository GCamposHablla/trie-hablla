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

func normalize(input string) string {
	if input == "" {
		return ""
	}

	builder := normBuilderPool.Get().(*strings.Builder)
	builder.Reset()
	defer normBuilderPool.Put(builder)

	builder.Grow(len(input))

	normalized := norm.NFD.String(input)

	var last rune
	for _, r := range normalized {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		if mapped, ok := leetMap[r]; ok {
			r = mapped
		}
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			r = unicode.ToLower(r)
			if r != last {
				builder.WriteRune(r)
				last = r
			}
			// se igual, ignora → colapso de repetições
		}
	}

	return builder.String()
}