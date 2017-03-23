package trie

import (
	"bufio"
	"os"
	"testing"
	"unicode/utf8"
)

func TestNewTrieNode(t *testing.T) {
	for r := rune(0); r < utf8.MaxRune; r++ {
		clean := int32(r)%2 == 0
		tn := newTrieNode(r, clean)

		if tn.val != r {
			t.Errorf("trieNode.val Error: Expected: %v, Actual: %v", r, tn.val)
		}
		if tn.clean != clean {
			t.Errorf("trieNode.clean Error: Expected: %b, Actual: %b", clean, tn.clean)
		}
		if tn.children == nil {
			t.Error("trieNode.children Error: Should not be nil.")
		}
		if len(tn.children) != 0 {
			t.Errorf("trieNode.children Error: Children map should be empty, actual value is %d", len(tn.children))
		}
	}
}

func TestNewTrieTree(t *testing.T) {
	endRunes := []rune{' ', '*', 'a', '$', '\n', 'r'}

	for _, er := range endRunes {
		tt := NewTrieTree(er)

		if tt.specialEndRune != er {
			t.Errorf("TrieTree.specialEndRune Error: Expected: %v, Actual: %v", er, tt.specialEndRune)
		}

		if tt.root.val != er {
			t.Errorf("TrieTree.root.val Error: Expected: %v, Actual: %v", er, tt.root.val)
		}
		if !tt.root.clean {
			t.Error("TrieTree.root.clean Error: Root should always be clean")
		}
		if tt.root.children == nil {
			t.Error("TrieTree.trieNode.children Error: should not be nil.")
		}
		if len(tt.root.children) != 0 {
			t.Errorf("TrieTree.trieNode.children Error: Children map should be empty, actual value is %d", len(tt.root.children))
		}
	}
}

func TestAdd(t *testing.T) {
	tt := NewTrieTree('*')

	if !tt.Add("") {
		t.Error("Add returned false for word with 0 length")
	}

	words := []string{"Robert", "Bad", "wold", "abcdefghijklmnopqrstuvwxyz"}
	for _, word := range words {
		if !tt.Add(word) {
			t.Errorf("Failed to add word %s to TrieTree", word)
		}

		// Now we need to check that it was added
		n := tt.root
		for c := 0; c < len(word); c++ {
			r := rune(word[c])
			if _, ok := n.children[r]; !ok {
				t.Error("Char %v was not added to tree", c)
			}

			n = n.children[r]
		}

		if _, ok := n.children['*']; !ok {
			t.Error("SpecialEndRune was not put at the end of the word")
		}
	}

	// Check that we get false on re-add of words
	for _, word := range words {
		if tt.Add(word) {
			t.Errorf("Word %s was already added, should get false back", word)
		}
	}

	// Check that word with special end rune returns false
	if tt.Add("test*word") {
		t.Error("Should get false back for words with special end rune in it")
	}
	// Now check that the tree is still clean
	if !checkTrieNodeIsClean(tt.root, tt.specialEndRune) {
		t.Error("Tree has been left in a dirty state after an invalid word added")
	}

	// Check that word with special end rune returns false
	tt.Add("test")
	if tt.Add("test*") {
		t.Error("Should get false back for words with special end rune in it")
	}
	// Now check that the tree is still clean
	if !checkTrieNodeIsClean(tt.root, tt.specialEndRune) {
		t.Error("Tree has been left in a dirty state after an invalid word added")
	}
}

func TestContains(t *testing.T) {
	tt := NewTrieTree('*')

	words := []string{"Robert", "Tardis", "testing", "babaloo", "golang"}

	for idx, word := range words {
		tt.Add(word)

		for jdx, testWord := range words {
			if jdx <= idx && !tt.Contains(testWord) {
				t.Errorf("Tree should contain word %s", testWord)
				return
			} else if jdx > idx && tt.Contains(testWord) {
				t.Errorf("Tree should not contain word %s", testWord)
				return
			}
		}
	}

	fo, err := os.Open("testdata/dict.txt")
	defer fo.Close()
	if err != nil {
		t.Error(err)
		return
	}

	scanner := bufio.NewScanner(fo)
	var dictWords []string
	for scanner.Scan() {
		dictWords = append(dictWords, scanner.Text())
	}

	bigt := NewTrieTree('*')
	for idx, word := range dictWords {
		bigt.Add(word)

		for jdx, testWord := range dictWords {
			if jdx <= idx && !bigt.Contains(testWord) {
				t.Errorf("Tree should contain word %s", testWord)
				return
			} else if jdx > idx && bigt.Contains(testWord) {
				t.Errorf("Tree should not contain word %s", testWord)
				return
			}
		}
	}
}

func checkTrieNodeIsClean(tn *trieNode, specialEndRune rune) bool {
	for cr, ctn := range tn.children {
		if cr == specialEndRune {
			if ctn != nil {
				return false
			}
			continue
		}

		if !ctn.clean {
			return false
		}

		if !checkTrieNodeIsClean(ctn, specialEndRune) {
			return false
		}
	}

	return true
}
