package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// Node represents an element in the graph
type Node struct {
	name     string
	children []string
}

func main() {
	filename := "example.txt" // Change this to the name of the file you want to read

	content, err := readFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	lines := splitStringByNewline(string(*content))
	nodes := convertSliceToNodes(processLinesSlice(lines))
	visualizeGraph(nodes)
}

func readFile(filename string) (*[]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func splitStringByNewline(input string) *[]string {
	lines := strings.Split(input, "\n")
	return &lines
}

func processLinesSlice(input *[]string) *[][]string {
	var result [][]string
	var currentGroup []string

	for _, line := range *input {
		if line == "" {
			if len(currentGroup) > 0 {
				result = append(result, currentGroup)
				currentGroup = nil
			}
		} else {
			currentGroup = append(currentGroup, line)
		}
	}

	if len(currentGroup) > 0 {
		result = append(result, currentGroup)
	}

	return &result
}

func convertSliceToNodes(input *[][]string) *[]Node {
	var result []Node
	for _, nodes := range *input {
		result = append(result, Node{name: nodes[0], children: nodes[1:]})
	}

	return &result
}

func findParentNodes(nodes *[]Node) *[]string {
	childNodes := make(map[string]bool)

	for _, node := range *nodes {
		for _, name := range node.children {
			childNodes[name] = true
		}
	}

	parentNodes := []string{}
	for _, parent := range *nodes {
		if !childNodes[parent.name] {
			parentNodes = append(parentNodes, parent.name)
		}
	}

	return &parentNodes
}

func findTopNodes(parent string, nodes *[]Node) []string {
	if parent == "" {
		return *findParentNodes(nodes)
	}
	for _, node := range *nodes {
		if node.name == parent {
			return node.children
		}
	}
	return nil
}

func generatePrefix(depth int, isLast []bool) string {
	prefix := ""
	for i := 0; i < depth; i++ {
		if isLast[i] {
			prefix += "    "
		} else {
			prefix += "│   "
		}
	}
	return prefix
}

func getBranch(index int, size int) string {
	if size == index+1 {
		return "└──"
	}
	return "├──"
}

func removeNode(nodes *[]Node, parent string) *[]Node {
	filteredNodes := []Node{}
	for _, node := range *nodes {
		if parent != node.name {
			filteredNodes = append(filteredNodes, node)
		}
	}
	return &filteredNodes
}

func visualizeParents(depth int, parentNodes []string, nodes *[]Node, isLast []bool) {
	if len(parentNodes) == 0 {
		return
	}
	currentIndex := 0
	for _, parent := range parentNodes {
		isLast := append(isLast, currentIndex == len(parentNodes)-1)
		prefix := generatePrefix(depth, isLast)
		branch := getBranch(currentIndex, len(parentNodes))
		fmt.Printf("%s%s %s\n", prefix, branch, parent)
		nextParents := findTopNodes(parent, nodes)
		nextNodes := removeNode(nodes, parent)
		visualizeParents(depth+1, nextParents, nextNodes, isLast)
		currentIndex++
	}
}

func visualizeGraph(nodes *[]Node) {
	depth := 0
	var parent string
	parentNodes := findTopNodes(parent, nodes)
	isLast := []bool{}
	fmt.Println(".")

	visualizeParents(depth, parentNodes, nodes, isLast)
}
