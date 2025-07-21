package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var proibidas = map[string]bool{
	"idiota":  true,
	"burro":   true,
	"nazismo": true,
}

func censurarMensagem(msg string) string {
	palavras := strings.Fields(msg)
	for i, p := range palavras {
		// Normaliza a palavra para comparar (caso queira ignorar maiúsculas)
		normalizada := strings.ToLower(p)
		if proibidas[normalizada] {
			palavras[i] = strings.Repeat("*", len(p)) // ou palavras[i] = "#proibido#"
		}
	}
	return strings.Join(palavras, " ")
}

func main() {
	var msg string
	fmt.Print("Digite sua mensagem: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		msg = scanner.Text()
	}

	msgCensurada := censurarMensagem(msg)
	fmt.Println("Mensagem filtrada:", msgCensurada)
}
