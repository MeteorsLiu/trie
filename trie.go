package trie

import (
	"fmt"
	"strings"
	"sync"
)

type Trie[T any] struct {
	sync.RWMutex
	root *Node[T]
}

func NewTrie[T any]() *Trie[T] {
	return &Trie[T]{root: NewNodeNil[T](nil)}
}

func (t *Trie[T]) Insert(k string, value T) {
	parts := strings.Split(k, ".")

	t.Lock()
	defer t.Unlock()

	if parts[0] == "+" {
		t.insert(parts[1:], value)
		parts[0] = ""
		t.insert(parts, value)
	} else {
		t.insert(parts, value)
	}
}

func (t *Trie[T]) insert(parts []string, value T) {
	node := t.root
	for i := len(parts) - 1; i >= 0; i-- {
		node = node.GetOrSet(parts[i], NewNodeNil[T](node))
	}
	node.MarkAsLeaf()
	node.SetData(value)
}

func (t *Trie[T]) search(node *Node[T], parts []string) *Node[T] {
	if len(parts) == 0 {
		return node
	}

	if c := node.Get(parts[len(parts)-1]); c != nil {
		if n := t.search(c, parts[:len(parts)-1]); n != nil && n.IsLeaf() {
			return n
		}
	}

	if c := node.GetWildcard(); c != nil {
		if n := t.search(c, parts[:len(parts)-1]); n != nil && n.IsLeaf() {
			return n
		}
	}

	return nil
}

func (t *Trie[T]) Search(k string) *Node[T] {
	parts := strings.Split(k, ".")

	t.RLock()
	n := t.search(t.root, parts)
	t.RUnlock()

	return n
}

func (t *Trie[T]) Remove(k string) {
	parts := strings.Split(k, ".")

	t.Lock()
	defer t.Unlock()

	n := t.search(t.root, parts)

	for i := 0; i < len(parts) && n != nil && n.IsLeaf(); i++ {
		n.Remove(parts[i])
		if n.parent != nil && len(n.children) < 2 {
			n.parent.RemoveNode(n)
		}
		n = n.parent
	}

}

// 后序历遍
func (t *Trie[T]) walk(node *Node[T], depth int, f func(int, string, *Node[T])) {
	if node == nil {
		return
	}
	node.ForEach(func(s string, n *Node[T]) {
		if !node.IsEmpty() {
			t.walk(n, depth+1, f)
		}
		f(depth, s, n)
	})
}

func (t *Trie[T]) Walk(f func(int, string, *Node[T])) {
	t.walk(t.root, 0, f)
}

func (t *Trie[T]) print(n *Node[T], key string, space int) {
	if n == nil {
		return
	}
	space += 10

	n.ForEach(func(s string, nc *Node[T]) {
		t.print(nc, s, space)
	})

	for i := 0; i < space; i++ {
		fmt.Print(" ")
	}

	if key == "" {
		if n != t.root {
			fmt.Printf("+: %v\n", n.Data)
		} else {
			fmt.Print("root\n")
		}
	} else {
		fmt.Printf("%s: %v\n", key, n.Data)
	}

}

func (t *Trie[T]) Print() {
	t.print(t.root, "", 0)
}
