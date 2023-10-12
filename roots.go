package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Node represents an element in the graph
type Node struct {
	name     string
	children []string
}

func main() {
	var progname string = os.Args[0]
	var input []string = os.Args[1:]

	flags, usage := parseFlags(progname, input)
	if len(input) == 0 {
		usage()
		stop("", 1)
	}
	if flags.help {
		usage()
		stop("", 0)
	}
	if flags.version {
		stop(version, 0)
	}

	fileName := input[0]
	content, err := readFile(fileName)
	check(err)

	lines := splitStringByNewline(string(*content))
	nodes := convertSliceToNodes(processLinesSlice(lines))
	validateNodes(nodes)
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

func findRootNodeNames(nodes *[]Node) *[]string {
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

func findChildrenNodeNames(nodes *[]Node, parent string) *[]string {
	if parent == "" {
		return findRootNodeNames(nodes)
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

func validateNodes(nodes *[]Node) {
	warningCount := 0
	warningCount += checkDuplicateParent(nodes)
	warningCount += checkDuplicateChild(nodes)
	if warningCount > 0 {
		warn(fmt.Sprintf("\nFound %d warning(s).", warningCount))
	}
}

func checkDuplicateParent(nodes *[]Node) int {
	uniqueNodes := []string{}
	warningCount := 0
	for _, node := range *nodes {
		if !contains(uniqueNodes, node.name) {
			uniqueNodes = append(uniqueNodes, node.name)
		} else {
			warn(fmt.Sprintf("Duplicate parent \"%s\".", node.name))
			warningCount++
		}
	}
	return warningCount
}

func checkDuplicateChild(nodes *[]Node) int {
	warningCount := 0
	for _, node := range *nodes {
		uniqueChildren := []string{}
		for _, child := range node.children {
			if !contains(uniqueChildren, child) {
				uniqueChildren = append(uniqueChildren, child)
			} else {
				warn(fmt.Sprintf("Duplicate child \"%s\" in node \"%s\".", child, node.name))
				warningCount++
			}
		}
	}
	return warningCount
}

func visualizeGraphLevel(nodes *[]Node, levelNodeNames *[]string, isLastNodes *[]bool) {
	for index, parent := range *levelNodeNames {
		isLast := index == len(*levelNodeNames)-1
		isLastNodes := append(*isLastNodes, isLast)
		prefix := generatePrefix(&isLastNodes)
		branch := getBranch(isLast)
		fmt.Printf("%s%s %s\n", prefix, branch, parent)

		childrenNodeNames := findChildrenNodeNames(nodes, parent)
		remainingNodes := removeNode(nodes, parent)
		visualizeGraphLevel(remainingNodes, childrenNodeNames, &isLastNodes)
	}
}

func visualizeGraph(nodes *[]Node) {
	var parent string
	rootLevelNodeNames := findChildrenNodeNames(nodes, parent)
	isLastNodes := []bool{}
	fmt.Println(".")

	visualizeGraphLevel(nodes, rootLevelNodeNames, &isLastNodes)
}
