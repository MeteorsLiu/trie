package trie

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

var localIP = net.IP{127, 0, 0, 1}

func TestTrie_Basic(t *testing.T) {

	tree := NewTrie[net.IP]()
	domains := []string{
		"example.com",
		"google.com",
		"localhost",
	}

	for _, domain := range domains {
		tree.Insert(domain, localIP)
	}

	node := tree.Search("example.com")
	assert.NotNil(t, node)
	assert.True(t, node.Data.Equal(localIP))
	//assert.NotNil(t, tree.Insert("", localIP))
	assert.Nil(t, tree.Search(""))
	assert.NotNil(t, tree.Search("localhost"))
	assert.Nil(t, tree.Search("www.google.com"))

	//t.Log(tree.Search("www.google.com"))
}

func TestTrie_RemoveBasic(t *testing.T) {

	tree := NewTrie[net.IP]()
	domains := []string{
		"example.com",
		"google.com",
		"localhost",
		"www.localhost.com",
	}

	for _, domain := range domains {
		tree.Insert(domain, localIP)
	}
	tree.Remove("example.com")
	tree.Walk(func(depth int, s string, n *Node[net.IP]) {
		t.Log(depth, s, n)
	})
	assert.Nil(t, tree.Search("example.com"))
	assert.NotNil(t, tree.Search("localhost"))
	assert.NotNil(t, tree.Search("google.com"))

	tree.Remove("google.com")
	assert.NotNil(t, tree.Search("localhost"))
	assert.Nil(t, tree.Search("google.com"))
	tree.Remove("localhost")
	assert.Nil(t, tree.Search("localhost"))
	t.Log(tree.Search("www.google.com"))
}

func TestTrie_Wildcard(t *testing.T) {
	tree := NewTrie[net.IP]()

	domains := []string{
		"*.example.com",
		"sub.*.example.com",
		"*.dev",
		".org",
		".example.net",
		".apple.*",
		"+.foo.com",
		"+.stun.*.*",
		"+.stun.*.*.*",
		"+.stun.*.*.*.*",
		"stun.l.google.com",
	}

	for _, domain := range domains {
		tree.Insert(domain, localIP)
	}
	tree.Walk(func(depth int, s string, n *Node[net.IP]) {
		t.Log(depth, s)
	})
	assert.NotNil(t, tree.Search("sub.example.com"))
	assert.NotNil(t, tree.Search("sub.foo.example.com"))
	assert.NotNil(t, tree.Search("test.org"))
	assert.NotNil(t, tree.Search("test.example.net"))
	assert.NotNil(t, tree.Search("test.apple.com"))
	assert.NotNil(t, tree.Search("test.foo.com"))
	assert.NotNil(t, tree.Search("foo.com"))
	assert.NotNil(t, tree.Search("global.stun.website.com"))
	assert.Nil(t, tree.Search("foo.sub.example.com"))
	assert.Nil(t, tree.Search("foo.example.dev"))
	assert.Nil(t, tree.Search("example.com"))
}

func TestTrie_RemoveWildcard(t *testing.T) {
	tree := NewTrie[net.IP]()

	domains := []string{
		"*.example.com",
		"sub.*.example.com",
		"*.dev",
		".org",
		".example.net",
		".apple.*",
		"+.foo.com",
		"+.stun.*.*",
		"+.stun.*.*.*",
		"+.stun.*.*.*.*",
		"stun.google.*.*",
	}

	for _, domain := range domains {
		tree.Insert(domain, localIP)
	}

	//tree.Print()

	assert.NotNil(t, tree.Search("sub.example.com"))
	assert.NotNil(t, tree.Search("sub.foo.example.com"))
	assert.NotNil(t, tree.Search("test.org"))
	assert.NotNil(t, tree.Search("test.example.net"))
	assert.NotNil(t, tree.Search("test.apple.com"))
	assert.NotNil(t, tree.Search("test.foo.com"))
	assert.NotNil(t, tree.Search("foo.com"))
	assert.NotNil(t, tree.Search("global.stun.website.com"))
	assert.Nil(t, tree.Search("foo.sub.example.com"))
	assert.Nil(t, tree.Search("foo.example.dev"))
	assert.Nil(t, tree.Search("example.com"))

	tree.Remove("*.example.com")
	assert.Nil(t, tree.Search("sub.example.com"))
	assert.Nil(t, tree.Search("sub.foo.example.com"))

	tree.Remove(".org")
	assert.Nil(t, tree.Search("test.org"))

	tree.Remove(".example.net")
	assert.Nil(t, tree.Search("test.example.net"))

	tree.Remove(".apple.*")
	assert.Nil(t, tree.Search("test.apple.com"))

	tree.Remove("+.foo.com")

	assert.Nil(t, tree.Search("test.foo.com"))
	assert.Nil(t, tree.Search("foo.com"))
	t.Log("start del stun")
	tree.Remove("+.stun.*.*.*.*")
	assert.NotNil(t, tree.Search("global.stun.website.com"))
	tree.Remove("+.stun.*.*.*")
	assert.NotNil(t, tree.Search("global.stun.website.com"))

	tree.Remove("+.stun.*.*")
	assert.Nil(t, tree.Search("global.stun.website.com"))
	assert.NotNil(t, tree.Search("stun.google.baidu.com"))
}

func TestTrie_Boundary(t *testing.T) {
	tree := NewTrie[net.IP]()
	tree.Insert("*.dev", localIP)

	tree.Insert(".", localIP)
	tree.Insert("..dev", localIP)
	assert.Nil(t, tree.Search("dev"))
}

func TestTrie_RemoveWildcard2(t *testing.T) {
	tree := NewTrie[net.IP]()

	domains := []string{
		"+.stun.*.*",
		"+.stun.*.*.*",
		"+.stun.*.*.*.*",
		"stun.l.google.com",
	}

	for _, domain := range domains {
		tree.Insert(domain, localIP)
	}

	tree.Remove("+.stun.*.*.*.*")
	assert.NotNil(t, tree.Search("global.stun.website.com"))

	tree.Remove("+.stun.*.*.*")
	assert.NotNil(t, tree.Search("global.stun.website.com"))
	tree.Print()
	tree.Remove("+.stun.*.*")
	tree.Print()
	assert.Nil(t, tree.Search("global.stun.website.com"))
}
