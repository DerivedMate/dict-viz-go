package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strings"
)

// A Node of a binary tree with entry weights
// The first node should
type Node struct {
	EntryAngle float64
	Word       string
	Parent     *Node
	Left       *Node
	Right      *Node
}

// NewRootNode : The constructor of the first node of a binary tree
func NewRootNode(word string) *Node {
	return &Node{
		EntryAngle: 0,
		Word:       word,
		Parent:     nil,
		Left:       nil,
		Right:      nil,
	}
}

// NewNode : The constructor of the Node type
func NewNode(parent *Node, word string) *Node {
	node := &Node{
		EntryAngle: 0.0,
		Word:       word,
		Parent:     parent,
		Left:       nil,
		Right:      nil,
	}

	node.EntryAngle = math.Acos(CosStringsAlphabet(parent.Word, word))

	return node
}

// GetUniqueSignsOfSnapshot : Finds all the unique characters of a snapshot to later use in vectorization of words. Kept for future use / debugging.
func GetUniqueSignsOfSnapshot(fileName string) string {
	f, err := os.Open(fileName)
	defer f.Close()
	HandleError(err)

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	HandleError(err)

	signs := ""

	for _, rec := range records[1:] {
		letters := strings.Split(rec[0], "")

		for _, l := range letters {
			l = strings.ToLower(l)
			if isRep := strings.Contains(signs, l); !isRep {
				signs += l
			}
		}
	}

	return signs
}

func MakeTree(root *Node, words []string) {
	n := int(math.Min(float64(len(words)), 2))
	for i := 0; i < len(words) && i < 2; i++ {
		n := NewNode(root, words[i])
		if root.Left != nil && root.Left.EntryAngle > n.EntryAngle {
			root.Right = root.Left
		}
		if root.Left != nil {
			root.Right = n
		} else {
			root.Left = n
		}
	}

	words = words[n:]
	m := len(words) / 2

	if root.Left != nil {
		MakeTree(root.Left, words[:m])
	}
	if root.Right != nil {
		MakeTree(root.Right, words[m:])
	}
}

func PrintTree(root *Node, i int) {
	if root == nil {
		return
	}

	fmt.Printf("[%v] {%v %v}\n", i, root.Word, math.Cos(root.EntryAngle))
	PrintTree(root.Left, i+1)
	PrintTree(root.Right, i+1)
}

func VectorizationMain() {
	f, err := os.Open("out-clean.csv")
	defer f.Close()
	HandleError(err)

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	HandleError(err)

	words := make([]string, 0)
	for _, r := range records[1:] {
		words = append(words, r[0])
	}

	root := NewRootNode(words[0])
	words = words[1:]
	MakeTree(root, words)
	PrintTree(root, 0)
}
