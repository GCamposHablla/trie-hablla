package main

import (
	"testing"
)

var testTexts = []string{
	"Este é um texto normal sem palavras ruins",
	"Texto com palavra4 pr0ib1da e outras c0isas",
	"Um texto muito longo com várias palavras que podem ou não ser censuradas",
	"palavra censur4da aqui",
	"",
	"a",
}

func BenchmarkFilter(b *testing.B) {
	censor, err := NewCensor("palavras.txt")
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		text := testTexts[i%len(testTexts)]
		censor.Filter(text)
	}
}

func BenchmarkNormalization(b *testing.B) {
	text := "T3xt0 c0m l33t sp34k e acentuação çãõ"
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		normalize(text)
	}
}