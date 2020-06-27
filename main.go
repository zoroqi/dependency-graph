package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	graphStr := graph()
	parseGraph(graphStr)
}

type pkg struct {
	name         string
	ver          string
	incompatible bool
	dep          []*pkg
}

func (m pkg) String() string {
	return m.name + "@" + m.ver
}

func parseGraph(graph string) {
	lines := strings.Split(graph, "\n")
	root := findRoot(lines)
	if root == "" {
		fmt.Println("no dependency")
		return
	}

	depMapping := make(map[string]*pkg)

	for _, line := range lines {
		if "" == line {
			continue
		}
		ss := strings.SplitN(line, " ", 2)
		lib := parsePkg(ss[0])
		if l, exist := depMapping[lib.String()]; exist {
			lib = l
		} else {
			depMapping[lib.String()] = lib
		}
		depModel := parsePkg(ss[1])
		if l, exist := depMapping[depModel.String()]; exist {
			depModel = l
		} else {
			depMapping[depModel.String()] = depModel
		}
		lib.dep = append(lib.dep, depModel)
	}

	sb := toString(depMapping[root])
	fmt.Println(sb.String())
}

func findRoot(lines []string) string {
	for _, l := range lines {
		if l != "" {
			return strings.Split(l, " ")[0] + "@"
		}
	}
	return ""
}

func parsePkg(str string) *pkg {
	lv := strings.SplitN(str, "@", 2)
	if len(lv) == 1 {
		lv = append(lv, "")
	}
	v := strings.SplitN(lv[1], "+", 2)
	depModel := pkg{name: lv[0],
		ver:          v[0],
		incompatible: len(v) == 2,
		dep:          make([]*pkg, 0),
	}
	return &depModel
}

// build space
func levelStr(level int) string {
	return strings.Repeat("    |", level)
}

type stackNode struct {
	pkg   *pkg
	index int // index in traversal
}

func toString(pkg *pkg) *strings.Builder {
	sb := strings.Builder{}
	stack := make([]*stackNode, 0, 10)

	stackMap := make(map[string]bool)
	push := func(l *stackNode) {
		stack = append(stack, l)
		stackMap[l.pkg.String()] = true
	}

	pop := func() *stackNode {
		if len(stack) < 0 {
			return nil
		}
		r := stack[len(stack)-1]
		stackMap[r.pkg.String()] = false
		stack = stack[0 : len(stack)-1]
		return r
	}

	top := func() *stackNode {
		if len(stack) == 0 {
			return nil
		}
		return stack[len(stack)-1]
	}

	length := func() int {
		return len(stack)
	}

	push(&stackNode{
		pkg:   pkg,
		index: 0,
	})

	push2 := func(tmp *stackNode) {
		if stackMap[tmp.pkg.dep[tmp.index].String()] {
			fmt.Printf("%s-%s:circular\n", levelStr(length()), tmp.pkg.dep[tmp.index].String())
		} else {
			push(&stackNode{
				pkg:   tmp.pkg.dep[tmp.index],
				index: 0,
			})
		}
		tmp.index++
	}

	for tmp := top(); tmp != nil; tmp = top() {
		if tmp.index == 0 {
			fmt.Printf("%s-%s\n", levelStr(length()-1), tmp.pkg.String())
			if len(tmp.pkg.dep) == 0 {
				pop()
			} else {
				push2(tmp)
			}
		} else {
			if tmp.index < len(tmp.pkg.dep) {
				push2(tmp)
			} else {
				pop()
			}
		}
	}
	return &sb
}

// execute go mod graph
func graph() string {
	if _, err := os.Stat("./go.mod"); os.IsNotExist(err) {
		fmt.Println("cannot find go.mod")
		os.Exit(1)
	}
	cmd := exec.Command("go", "mod", "graph")
	resultBytes, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(resultBytes)
}

func list() string {
	if _, err := os.Stat("./go.mod"); os.IsNotExist(err) {
		fmt.Println("cannot find go.mod")
		os.Exit(1)
	}
	cmd := exec.Command("go", "list", "-m", "all")
	resultBytes, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(resultBytes)
}
