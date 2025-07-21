package main

import (
	"fmt"
	"strings"
	"unicode"
	"bufio"
	"os"
	"golang.org/x/text/unicode/norm"
)

var leetMap = map[rune]rune{
    '4': 'a', '@': 'a', '^': 'a',                       // A
    '8': 'b',                                          // B
    '<': 'c', '(': 'c',                               // C
    ')': 'd',                              			 // D
    '3': 'e', '&': 'e', '€': 'e',                     // E
    'ƒ': 'f',                                         // F
    '6': 'g',                                          // G
    '#': 'h',                                          // H
    '1': 'i', '!': 'i',                               // I
    '£': 'l',                              			 // L
    '0': 'o',                                          // O
    '5': 's', '$': 's',                               // S
    '7': 't', '+' : 't',                              // T
    'µ': 'u',                                          // U
    '2': 'z',                                          // Z
}


// Nó da Trie
type TrieNode struct {
	letter   rune
	children map[rune]*TrieNode
	isEnd    bool
	self     *TrieNode 
}

// Trie com nós que não repetem letras consecutivas
type Trie struct {
	root *TrieNode
}

// Cria novo nó
func newTrieNode(letter rune) *TrieNode {
	node := &TrieNode{
		letter:   letter,
		children: make(map[rune]*TrieNode),
		isEnd:    false,
	}
	node.self = node // Aponta para si mesmo
	return node
}

// Cria nova trie
func NewTrie() *Trie {
	return &Trie{root: newTrieNode(0)} // letra 0 para raiz
}

// Insere palavra, evitando letras consecutivas repetidas
func (t *Trie) Insert(word string) {
	node := t.root
	var lastLetter rune

	for i, ch := range word {
		if i > 0 && ch == lastLetter {
			// Letra repetida: usar self (não cria novo filho)
			node = node.self
			continue
		}
		lastLetter = ch

		if _, exists := node.children[ch]; !exists {
			node.children[ch] = newTrieNode(ch)
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

// Busca com uso de self se letra se repete
func (t *Trie) Search(phrase string) string {
	words := strings.Split(phrase,` `)
	fmt.Println(words)
	for i, word := range words {
		normalized := normalizeInput(word)
		found := t.searchRec(t.root, []rune(normalized), "", 0)
		if found {
			words[i] = strings.Repeat("*", len(word))
		}
	}
	return strings.Join(words, " ")
}

func (t *Trie) searchRec(node *TrieNode, word []rune, caminho string, pos int) bool {
	if len(word) == 0 {
		fmt.Printf("🚩 Caminho final: %s\n", caminho)
		return node.isEnd
	}

	ch := word[0]

	// Caso 1: letra está nos filhos
	if child, ok := node.children[ch]; ok {
		fmt.Printf("➡️  Avançando: %s%c (via filho)\n", caminho, ch)
		return t.searchRec(child, word[1:], caminho+string(ch), pos+1)
	}

	// Caso 2: letra igual à atual e pode usar self
	if node.letter == ch && node.self != nil {
		fmt.Printf("🔁 Usando self: %s%c (reutilizando '%c')\n", caminho, ch, node.letter)
		return t.searchRec(node.self, word[1:], caminho+string(ch), pos+1)
	}

	fmt.Printf("❌ Falhou em: %s%c\n", caminho, ch)
	return false
}

func normalizeInput(input string) string {
	// 1. Normaliza para decompor acentos (NFD)
	normStr := norm.NFD.String(input)
	var builder strings.Builder

	for _, r := range normStr {
		// 2. Remove marcas de acento
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		// 3. Substitui se estiver no leetMap
		if mapped, ok := leetMap[r]; ok {
			builder.WriteRune(mapped)
		} else if unicode.IsLetter(r) {
			builder.WriteRune(unicode.ToLower(r))
		}
		// Caracteres ignorados (pontuação, etc)
	}

	return builder.String()
}

// Lê um .txt e insere cada palavra na árvore
func (t *Trie) LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineno := 0
	for scanner.Scan() {
		lineno++
		raw := scanner.Text()
		word := normalizeInput(raw)
		if word == "" {
			continue // ignora linhas vazias ou com símbolos não convertidos
		}
		t.Insert(word)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("erro ao ler o arquivo: %v", err)
	}

	return nil
}

// Função principal
func main() {
	tree := NewTrie()
	reader := bufio.NewReader(os.Stdin)

	err := tree.LoadFromFile("palavras.txt")
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	var palavra string
	fmt.Print("Digite uma palavra: ")
	palavra, err2 := reader.ReadString('\n')
	if err2 != nil {
		fmt.Println("Erro ao ler entrada:", err2)
		return
	}

	fmt.Println("Resultado:", tree.Search(palavra))
}

