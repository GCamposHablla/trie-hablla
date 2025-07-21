package main

import (
	"fmt"
	"strings"
)

var palavrasProibidas = map[string]bool{
	"palavrão1": true,
	"palavrão2": true,
	"idiota":    true,
	"nazismo":   true,
}

func contemPalavraProibida(mensagem string) bool {
	palavras := strings.Fields(strings.ToLower(mensagem)) // separa por espaço
	for _, palavra := range palavras {
		if palavrasProibidas[palavra] {
			return true
		}
	}
	return false
}

func main() {
	var msg string
	fmt.Print("Digite a mensagem: ")
	fmt.Scanln(&msg)

	if contemPalavraProibida(msg) {
		fmt.Println("🚫 Mensagem contém palavras proibidas!")
	} else {
		fmt.Println("✅ Mensagem liberada.")
	}
}
