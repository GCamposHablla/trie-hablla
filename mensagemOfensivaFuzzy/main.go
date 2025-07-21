package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/agnivade/levenshtein"
)

// Lista simples de palavras proibidas
var palavrasProibidas = map[string]bool{
	"idiota":  true,
	"burro":   true,
	"otário":  true,
	"nazismo": true,
	"caralho": true,
}

// Remove pontuação da palavra, deixa só letras e números
func removerPontuacao(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return -1
	}, s)
}

// Limita repetições consecutivas de caracteres a no máximo 2 (ex: "burrroooo" → "burroo")
func normalizarRepeticoes(s string) string {
	var resultado []rune
	var ultimo rune
	contagem := 0

	for _, r := range s {
		if r == ultimo {
			contagem++
			if contagem < 3 { // aceita até 2 repetições consecutivas
				resultado = append(resultado, r)
			}
		} else {
			ultimo = r
			contagem = 1
			resultado = append(resultado, r)
		}
	}
	return string(resultado)
}

// Função que usa distância de Levenshtein para detectar palavras ofensivas próximas
func mensagemOfensivaFuzzy(mensagem string) bool {
	palavras := strings.Fields(strings.ToLower(mensagem))

	for _, palavra := range palavras {
		palavra = removerPontuacao(palavra)
		palavra = normalizarRepeticoes(palavra)

		for proibida := range palavrasProibidas {
			distancia := levenshtein.ComputeDistance(palavra, proibida)
			if distancia <= 2 {
				fmt.Printf("Detectado parecido com palavra proibida: '%s' (distância %d)\n", palavra, distancia)
				return true
			}
		}
	}
	return false
}

func main() {
	testes := []string{
		"Você é um idiota!",
		"Ei, seu burrrooooo",
		"Olá, tudo bem?",
		"Otarião da Silva",
		"Que pessoa legal",
		"Nazista",
		"carlhooooo",
	}

	for _, msg := range testes {
		if mensagemOfensivaFuzzy(msg) {
			fmt.Printf("Mensagem bloqueada: %q\n", msg)
		} else {
			fmt.Printf("Mensagem liberada: %q\n", msg)
		}
	}
}
