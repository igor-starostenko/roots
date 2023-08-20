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

func findTopNodes(nodes *[]Node, parent string) *[]string {
	if parent == "" {
		return findParentNodes(nodes)
	}
	for _, node := range *nodes {
		if node.name == parent {
			return &(node.children)
		}
	}
	return &[]string{}
}

func generatePrefix(isLastNodes *[]bool) string {
	depth := len(*isLastNodes) - 1
	prefix := ""
	for i := 0; i < depth; i++ {
		if (*isLastNodes)[i] {
			prefix += "    "
		} else {
			prefix += "│   "
		}
	}
	return prefix
}

func getBranch(isLast bool) string {
	if isLast {
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

func visualizeGraphLevel(nodes *[]Node, levelNodeNames *[]string, isLastNodes *[]bool) {
	for index, parent := range *levelNodeNames {
		isLast := index == len(*levelNodeNames)-1
		isLastNodes := append(*isLastNodes, isLast)
		prefix := generatePrefix(&isLastNodes)
		branch := getBranch(isLast)
		fmt.Printf("%s%s %s\n", prefix, branch, parent)

		childrenNodes := findTopNodes(nodes, parent)
		remainingNodes := removeNode(nodes, parent)
		visualizeGraphLevel(remainingNodes, childrenNodes, &isLastNodes)
	}
}

func visualizeGraph(nodes *[]Node) {
	var parent string
	parentNodes := findTopNodes(nodes, parent)
	isLastNodes := []bool{}
	fmt.Println(".")

	visualizeGraphLevel(nodes, parentNodes, &isLastNodes)
}
