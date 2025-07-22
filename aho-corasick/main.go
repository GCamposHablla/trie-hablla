package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	censor, err := NewCensor("palavras.txt")
	if err != nil {
		fmt.Println("Erro ao inicializar censura:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Digite uma frase: ")
	text, _ := reader.ReadString('\n')

	found, censored := censor.Filter(text)
	if found {
		fmt.Println("Texto censurado:", censored)
	} else {
		fmt.Println("Nenhuma palavra sensível encontrada.")
	}
}
