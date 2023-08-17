package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	filename := "example.txt" // Change this to the name of the file you want to read

	content, err := readFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	lines := splitStringByNewline(string(content))
	nodes := processLinesSlice(lines)
	nodesMap := convertNodesToMap(nodes)
	visualizeGraph(nodesMap)

	// for key, values := range nodesMap {
	// 	fmt.Printf("Key: %s, Values: %v\n", key, values)
	// }
}

func readFile(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func splitStringByNewline(input string) []string {
	return strings.Split(input, "\n")
}

func processLinesSlice(input []string) [][]string {
	var result [][]string
	var currentGroup []string

	for _, line := range input {
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

	return result
}

func convertNodesToMap(input [][]string) map[string][]string {
	result := make(map[string][]string)

	for _, group := range input {
		if len(group) > 0 {
			key := group[0]
			values := group[1:]
			result[key] = values
		}
	}

	return result
}

func findParentNodes(nodesMap map[string][]string) []string {
	childNodes := make(map[string]bool)

	for _, children := range nodesMap {
		for _, node := range children {
			childNodes[node] = true
		}
	}

	parentNodes := []string{}
	for parent := range nodesMap {
		if !childNodes[parent] {
			parentNodes = append(parentNodes, parent)
		}
	}

	return parentNodes
}

func findTopNodes(parent string, nodesMap map[string][]string) []string {
	if parent == "" {
		return findParentNodes(nodesMap)
	}
	return nodesMap[parent]
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

func visualizeParents(depth int, parentNodes []string, nodes map[string][]string, isLast []bool) {
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
		delete(nodes, parent)
		visualizeParents(depth+1, nextParents, nodes, isLast)
		currentIndex++
	}
}

func visualizeGraph(nodes map[string][]string) {
	depth := 0
	var parent string
	parentNodes := findTopNodes(parent, nodes)
	isLast := []bool{}
	fmt.Println(".")

	visualizeParents(depth, parentNodes, nodes, isLast)
}
