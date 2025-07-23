package main

import (
	"testing"
)

func BenchmarkParallelFilter(b *testing.B) {
	censor, err := NewCensor("palavras.txt")
	if err != nil {
		b.Fatalf("erro inicializando censura: %v", err)
	}

	input := "essa frase contém naaaaazismo, mas não deveria"
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			found, _ := censor.Filter(input)
			_ = found // ou valide resultado se quiser
		}
	})
}