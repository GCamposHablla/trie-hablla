package main

import (
	"fmt"
	"strings"
)

// Nó da Radix Tree
type RadixNode struct {
	children map[string]*RadixNode
	isEnd    bool
}

// Árvore Radix
type RadixTree struct {
	root *RadixNode
}

// Novo nó
func newRadixNode() *RadixNode {
	return &RadixNode{
		children: make(map[string]*RadixNode),
		isEnd:    false,
	}
}

// Nova árvore
func NewRadixTree() *RadixTree {
	return &RadixTree{root: newRadixNode()}
}

// Inserção de palavra na Radix Tree
func (t *RadixTree) Insert(word string) {
	t.insertRec(t.root, word)
}

func (t *RadixTree) insertRec(node *RadixNode, word string) {
	for edge, child := range node.children {
		// Encontra o maior prefixo comum entre word e edge
		commonPrefix := getCommonPrefix(word, edge)
		if commonPrefix == "" {
			continue
		}

		if commonPrefix == edge {
			// Continua inserção no filho
			t.insertRec(child, word[len(commonPrefix):])
			return
		}

		// Quebra o nó existente
		newChild := newRadixNode()
		newChild.children[edge[len(commonPrefix):]] = child
		newChild.isEnd = child.isEnd

		// Atualiza filho original
		node.children[commonPrefix] = newChild
		delete(node.children, edge)

		child.isEnd = false

		// Insere parte restante da nova palavra, se houver
		if len(word[len(commonPrefix):]) > 0 {
			newNode := newRadixNode()
			newNode.isEnd = true
			newChild.children[word[len(commonPrefix):]] = newNode
		} else {
			newChild.isEnd = true
		}
		return
	}

	// Nenhum prefixo comum, adiciona como novo ramo
	newNode := newRadixNode()
	newNode.isEnd = true
	node.children[word] = newNode
}

// Busca uma palavra na árvore
func (t *RadixTree) Search(word string) bool {
	fmt.Printf("🔍 Buscando por: %s\n", word)
	return t.searchRec(t.root, word, "")
}

func (t *RadixTree) searchRec(node *RadixNode, word string, caminho string) bool {
	for edge, child := range node.children {
		if strings.HasPrefix(word, edge) {
			fmt.Printf("➡️  Caminho: %s%s\n", caminho, edge)
			if len(word) == len(edge) {
				return child.isEnd
			}
			return t.searchRec(child, word[len(edge):], caminho+edge)
		}
	}
	fmt.Printf("❌ Falhou em: %s\n", caminho+word)
	return false
}

// Função auxiliar para prefixo comum
func getCommonPrefix(a, b string) string {
	minLen := min(len(a), len(b))
	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return a[:i]
		}
	}
	return a[:minLen]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Exemplo de uso
func main() {
	tree := NewRadixTree()

	tree.Insert("banana")
	tree.Insert("bandido")
	tree.Insert("banda")
	tree.Insert("batata")

	fmt.Println("Resultado:", tree.Search("banana"))     // true
	fmt.Println("Resultado:", tree.Search("ban"))        // false
	fmt.Println("Resultado:", tree.Search("bandido"))    // true
	fmt.Println("Resultado:", tree.Search("batatinha"))  // false
}
