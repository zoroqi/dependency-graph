package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type exclude map[string]bool

func (e exclude) String() string {
	s := ""
	for k := range e {
		s += k
	}
	return s
}

func (e exclude) Set(s string) error {
	e[s] = true
	return nil
}

func main() {
	p := flag.String("p", "", `default: tree print
rl: reverse line print
rt: reverse tree print
wt: whole tree print
`)
	s := flag.String("s", "", "search pkg name")
	level := flag.Int("l", 0, "max level")
	var exPkg exclude
	exPkg = make(map[string]bool)
	flag.Var(&exPkg, "ex", "exclude package")
	var exPre exclude
	exPre = make(map[string]bool)
	flag.Var(&exPre, "expre", "exclude package, prefix match")

	flag.Parse()
	graphStr := graph()
	root := parseGraph(graphStr)
	if root == nil {
		fmt.Println("no dependency")
		return
	}
	tree := newTree(root)

	match := compoundedMatch(buildMath(*s, *level, exPkg, exPre)...)

	var sh stringHandler
	switch *p {
	case "rt":
		sh = reverseLevelString
	case "rl":
		sh = reverseLineString
	case "wt":
		sh = wholeLevelString(compoundedMatch(buildMath("", *level, exPkg, exPre)...))
	default:
		sh = levelString
	}

	str := treeString(tree, 0, match, sh)
	fmt.Println(str)
}

func buildMath(s string, level int, exPkg exclude, exPre exclude) []filterHandler {
	matches := make([]filterHandler, 0)
	if strings.TrimSpace(s) != "" {
		matches = append(matches, searchPackage(strings.TrimSpace(s)))
	}
	if level > 0 {
		matches = append(matches, maxLevelFilter(level))
	}
	if len(exPkg) > 0 {
		matches = append(matches, excludePackage(exPkg))
	}
	if len(exPre) > 0 {
		matches = append(matches, excludePrefixPackage(exPre))
	}
	return matches
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

func parseGraph(graph string) *pkg {
	lines := strings.Split(graph, "\n")
	root := findRoot(lines)
	if root == "" {
		return nil
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

	return depMapping[root]
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
