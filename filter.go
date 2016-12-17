package trie

import (
	"sync"
)

type Trie struct {
	Root *TrieNode
	mu   sync.Mutex
}

type TrieNode struct {
	Children map[rune]*TrieNode
	End      bool
}

func NewTrie() *Trie {
	t := &Trie{}
	t.Root = NewTrieNode()

	return t
}

func NewTrieNode() *TrieNode {
	n := &TrieNode{}
	n.Children = make(map[rune]*TrieNode)
	n.End = false

	return n
}

func (t *Trie) Lock() {
	t.mu.Lock()
}

func (t *Trie) Unlock() {
	t.mu.Unlock()
}

//ÐÂÔöÒª¹ýÂËµÄ´Ê
func (t *Trie) Add(txt string) {
	if len(txt) < 1 {
		return
	}
	chars := []rune(txt)
	slen := len(chars)
	node := t.Root
	for i := 0; i < slen; i++ {
		if _, exists := node.Children[chars[i]]; !exists {
			node.Children[chars[i]] = NewTrieNode()
		}
		node = node.Children[chars[i]]
	}
	node.End = true
}

//ÆÁ±Î×ÖËÑË÷²éÑ¯
func (t *Trie) Find(txt string) bool {
	chars := []rune(txt)
	slen := len(chars)
	node := t.Root
	for i := 0; i < slen; i++ {
		if _, exists := node.Children[chars[i]]; exists {
			node = node.Children[chars[i]]
			for j := i + 1; j < slen; j++ {
				if _, exists := node.Children[chars[j]]; !exists {
					break
				}
				node = node.Children[chars[j]]
				if node.End == true {
					return true
				}
			}
			node = t.Root
		}
	}
	return false
}

//ÆÁ±Î×ÖËÑË÷Ìæ»»
func (t *Trie) Replace(txt string) (string, []string) {
	chars := []rune(txt)
	result := []rune(txt)
	find := make([]string, 0, 10)
	slen := len(chars)
	node := t.Root
	for i := 0; i < slen; i++ {
		if _, exists := node.Children[chars[i]]; exists {
			node = node.Children[chars[i]]
			for j := i + 1; j < slen; j++ {
				if _, exists := node.Children[chars[j]]; !exists {
					break
				}
				node = node.Children[chars[j]]
				if node.End == true {
					for t := i; t <= j; t++ {
						result[t] = '*'
					}
					find = append(find, string(chars[i:j+1]))
					i = j
					node = t.Root
					break
				}
			}
			node = t.Root
		}
	}
	return string(result), find
}
