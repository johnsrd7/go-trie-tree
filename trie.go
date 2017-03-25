package trie

// trieNode is an internal node for the trie tree.
type trieNode struct {
	val      rune
	clean    bool
	children map[rune]*trieNode
}

// Tree is a trie tree ADT that holds words of an alphabet.
type Tree struct {
	root           *trieNode
	specialEndRune rune
}

// newTrieNode creates a new trieNode with the given value and
// the given clean bit set.
func newTrieNode(val rune, clean bool) *trieNode {
	return &trieNode{val, clean, make(map[rune]*trieNode)}
}

// NewTrieTree creates a new to a TrieTree that keeps the given
// specialEndRune as a special character that cannot be used in
// words given to the tree.
func NewTrieTree(specialEndRune rune) *Tree {
	return &Tree{newTrieNode(specialEndRune, true), specialEndRune}
}

// Add adds the given word to the tree.
func (t Tree) Add(word string) bool {
	if len(word) == 0 {
		return true
	}

	curNode := t.root
	undo := false
	for _, c := range word {
		if c == t.specialEndRune {
			// The word had the special end rune, so
			// we need to break here and then undo any
			// dirty nodes that were added.
			undo = true
			break
		}

		// First, we need to see if the char is in the map
		if _, ok := curNode.children[c]; !ok {
			curNode.children[c] = &trieNode{c, false, make(map[rune]*trieNode)}
		}

		// Next, set the curNode to the matching one in children
		// and "recurse" with it.
		curNode = curNode.children[c]
	}

	if undo {
		curNode := t.root
		for _, c := range word {
			if c == t.specialEndRune {
				if curNode.children[c] != nil {
					// Set it back to nil if it got messed up.
					curNode.children[c] = nil
				}

				break
			}

			if !curNode.children[c].clean {
				// We have found our dirty path,
				// so we can just delete that node
				// from the child and rid the whole
				// dirty path.
				delete(curNode.children, c)
				break
			}

			curNode = curNode.children[c]
		}

		return false
	}

	// At the end, we have the final node, so we need to set the
	// ending char so we know that we have a word.
	if _, ok := curNode.children[t.specialEndRune]; ok {
		// If we are here, then the word already exists,
		// so we don't need to do anything because the word
		// already exists.
		return false
	}

	// If we got here, then the path for this word doesn't already
	// exist, so we need to add the end rune.
	curNode.children[t.specialEndRune] = nil

	// At this point, we just need to set all the nodes to clean
	curNode = t.root
	for _, c := range word {
		curNode = curNode.children[c]
		curNode.clean = true
	}

	return true
}

// Contains returns true if the given word is contained in the tree.
func (t Tree) Contains(word string) bool {
	if len(word) == 0 {
		return true
	}

	curNode := t.root
	for _, c := range word {
		if _, ok := curNode.children[c]; !ok {
			return false
		}

		curNode = curNode.children[c]
	}

	// If we got here, then we got through the word, so we just
	// need to check if the specialEndRune is here. If not, then
	// the word hasn't been added yet.
	_, ok := curNode.children[t.specialEndRune]
	return ok
}

// Delete removes the given word from the tree.
func (t Tree) Delete(word string) {
	if len(word) == 0 {
		return
	}

	curNode := t.root
	for _, c := range word {
		if _, ok := curNode.children[c]; !ok {
			break
		}

		curNode = curNode.children[c]
	}

	delete(curNode.children, t.specialEndRune)
}
