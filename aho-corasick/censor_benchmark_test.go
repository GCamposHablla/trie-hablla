package main

import (
	"os"
	"strings"
	"testing"
)

var (
	censor      *Censor
	sampleInput = "Eu sou nazista vai se fuder caralho!"
)

func init() {
	var err error
	censor, err = NewCensor("palavras.txt")
	if err != nil {
		panic("Erro ao inicializar o censor no benchmark: " + err.Error())
	}
}

// Benchmark simples com input estático
func BenchmarkFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		found, _ := censor.Filter(sampleInput)
		if !found {
			b.Error("Texto de benchmark não foi censurado como esperado")
		}
	}
}

func BenchmarkFilterParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			censor.Filter(sampleInput)
		}
	})
}

func BenchmarkFilterMultipleLinesParallel(b *testing.B) {
	data, err := os.ReadFile("testdata/mensagens.txt")
	if err != nil {
		b.Fatalf("Erro ao ler arquivo de mensagens: %v", err)
	}

	// Quebra o texto por linhas (cada uma é uma mensagem independente)
	lines := strings.Split(string(data), "\n")
	if len(lines) < 100 {
		b.Fatalf("Poucas mensagens para paralelismo: %d", len(lines))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Roda em round-robin nas mensagens
			msg := lines[i%len(lines)]
			censor.Filter(msg)
			i++
		}
	})
}

// Benchmark com leitura de arquivo grande simulando múltiplas mensagens
func BenchmarkFilterMultipleLines(b *testing.B) {
	lines, err := os.ReadFile("testdata/mensagens.txt")
	if err != nil {
		b.Fatalf("Erro ao ler arquivo de mensagens: %v", err)
	}
	text := string(lines)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		censor.Filter(text)
	}
}
