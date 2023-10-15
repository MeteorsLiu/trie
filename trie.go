package trie

import (
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

func (t *Trie[T]) searchPath(node *Node[T], path *[]*Node[T], parts []string) *Node[T] {
	if len(parts) == 0 {
		*path = append(*path, node)
		return node
	}

	if c := node.Get(parts[len(parts)-1]); c != nil {
		*path = append(*path, node)
		if n := t.searchPath(c, path, parts[:len(parts)-1]); n != nil && n.IsLeaf() {
			return n
		}
	}

	if c := node.GetWildcard(); c != nil {
		*path = append(*path, node)
		if n := t.searchPath(c, path, parts[:len(parts)-1]); n != nil && n.IsLeaf() {
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
	var path []*Node[T]
	parts := strings.Split(k, ".")

	t.Lock()
	defer t.Unlock()
	t.searchPath(t.root, &path, parts)
	for i := len(path) - 1; i >= 0; i-- {
		each := path[i]
		if !each.IsLeaf() {
			next := path[i+1]
			if next.IsLeaf() {
				each.RemoveNode(next)
			}
		}
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
