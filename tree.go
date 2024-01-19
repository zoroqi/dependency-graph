package main

import (
	"fmt"
	"io"
	"strings"
)

type pkgTreeNode struct {
	name     string
	version  string
	parent   *pkgTreeNode
	dep      []*pkgTreeNode
	circular bool
	already  bool
}

func (m pkgTreeNode) String() string {
	return m.name + "@" + m.version
}

func (m pkgTreeNode) FlagString() string {
	var s string
	if m.circular {
		s += ":circular"
	}
	if m.already {
		s += ":already"
	}
	return s
}

// build space
func levelStr(level int) string {
	return strings.Repeat("    |", level)
}

type stackNode struct {
	pkg   *pkg
	index int // index in traversal
	node  *pkgTreeNode
}

func newTree(pkg *pkg) *pkgTreeNode {
	stack := make([]*stackNode, 0, 10)

	stackMap := make(map[string]bool)
	nodeMap := make(map[string]bool)

	push := func(l *stackNode) {
		stack = append(stack, l)
		stackMap[l.pkg.String()] = true
		nodeMap[l.pkg.String()] = true
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

	push(&stackNode{
		pkg:   pkg,
		index: 0,
		node:  newNode(pkg),
	})

	root := top().node

	push2 := func(tmp *stackNode) {
		n := newNode(tmp.pkg.dep[tmp.index])
		n.parent = tmp.node
		if stackMap[tmp.pkg.dep[tmp.index].String()] {
			n.circular = true
		} else if nodeMap[tmp.pkg.dep[tmp.index].String()] {
			n.already = true
		} else {
			push(&stackNode{
				pkg:   tmp.pkg.dep[tmp.index],
				index: 0,
				node:  n,
			})
		}
		tmp.node.dep = append(tmp.node.dep, n)
		tmp.index++
	}

	for tmp := top(); tmp != nil; tmp = top() {
		if tmp.index == 0 {
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
	return root
}

func newNode(pkg *pkg) *pkgTreeNode {
	return &pkgTreeNode{
		name:     pkg.name,
		version:  pkg.ver,
		parent:   nil,
		dep:      nil,
		circular: false,
		already:  false,
	}
}

type filterHandler func(level int, node *pkgTreeNode) (isPrint bool, isGoing bool)

type stringHandler func(level int, node *pkgTreeNode, sb io.StringWriter)

func maxLevelFilter(max int) filterHandler {
	return func(level int, node *pkgTreeNode) (bool, bool) {
		return true, level < max
	}
}

func excludePackage(p map[string]bool) filterHandler {
	return func(level int, node *pkgTreeNode) (bool, bool) {
		return true, !p[node.name]
	}
}

func excludePrefixPackage(p map[string]bool) filterHandler {
	return func(level int, node *pkgTreeNode) (bool, bool) {
		for k := range p {
			if strings.HasPrefix(node.name, k) {
				return true, false
			}
		}
		return true, true
	}
}

func excludeErrVersion(p []*pkg) filterHandler {
	m := make(map[string]bool)
	for _, r := range p {
		m[newNode(r).String()] = true
	}
	return func(level int, node *pkgTreeNode) (bool, bool) {
		if m[node.String()] {
			return true, true
		}
		return false, false
	}
}

func searchPackage(p string) filterHandler {
	return func(level int, node *pkgTreeNode) (bool, bool) {
		return p == node.name, true
	}
}

func compoundedMatch(filters ...filterHandler) filterHandler {
	return func(level int, node *pkgTreeNode) (bool, bool) {
		p, s := true, true
		for _, m := range filters {
			p1, s1 := m(level, node)
			p = p && p1
			s = s && s1
		}
		return p, s
	}
}

// - a
// - |-b
// -   |-c
func levelString(level int, n *pkgTreeNode, sb io.StringWriter) {
	sb.WriteString(fmt.Sprintf("%s-%s%s\n", levelStr(level), n.String(), n.FlagString()))
}

// - c
// - |-b
// -   |-a
func reverseLevelString(level int, n *pkgTreeNode, sb io.StringWriter) {
	p := n
	h := 0
	for p != nil {
		sb.WriteString(fmt.Sprintf("%s-%s%s", levelStr(h), p.String(), p.FlagString()))
		h++
		p = p.parent
		if p != nil {
			sb.WriteString("\n")
		}
	}
	sb.WriteString("\n")
}

// c -> b -> a
func reverseLineString(level int, n *pkgTreeNode, sb io.StringWriter) {
	p := n
	for p != nil {
		sb.WriteString(fmt.Sprintf("%s%s", p.String(), p.FlagString()))
		p = p.parent
		if p != nil {
			sb.WriteString(fmt.Sprintf(" -> "))
		}
	}
	sb.WriteString("\n")
}

// - a
// - |-b
// -   |-c
//
// - a
// - |-d
// -   |-e
func wholeLevelString(match filterHandler) stringHandler {
	return func(level int, node *pkgTreeNode, sb io.StringWriter) {
		sbParent := ""
		p := node.parent
		high := level - 1
		for p != nil {
			sbParent = fmt.Sprintf("%s-%s%s\n", levelStr(high), p.String(), p.FlagString()) + sbParent
			p = p.parent
			high--
		}
		sb.WriteString(sbParent)
		treeString(node, level, match, levelString, sb)
		sb.WriteString("\n")
	}
}

func treeString(root *pkgTreeNode, level int, match filterHandler, sh stringHandler, out io.StringWriter) {
	dfs(root, level, func(level int, node *pkgTreeNode) bool {
		p, g := match(level, node)
		if p {
			sh(level, node, out)
		}
		return g
	})
}

func dotString(actualDepend []*pkg) stringHandler {

	depend := make(map[string]bool)
	for _, r := range actualDepend {
		depend[newNode(r).String()] = true
	}

	increment := -1
	index := make(map[string]int)
	nodeStmt := func(node *pkgTreeNode, sb io.StringWriter) {
		i, exist := index[node.String()]
		if !exist {
			increment++
			i = increment
			index[node.String()] = i
			if depend[node.String()] {
				sb.WriteString(fmt.Sprintf("%d [label=\"%s\" style=\"filled\"]\n", i, node.String()))
			} else {
				sb.WriteString(fmt.Sprintf("%d [label=\"%s\"]\n", i, node.String()))
			}
		}
	}

	repeat := make(map[string]bool)
	return func(level int, node *pkgTreeNode, sb io.StringWriter) {
		tmp := node
		for tmp.parent != nil {
			nodeStmt(tmp, sb)
			nodeStmt(tmp.parent, sb)
			i := index[tmp.String()]
			pi := index[tmp.parent.String()]
			ss := fmt.Sprintf("%d -> %d;\n", pi, i)
			if !repeat[ss] {
				sb.WriteString(ss)
				repeat[ss] = true
			}
			tmp = tmp.parent
		}
	}
}

func mermaidString(actualDepend []*pkg) stringHandler {

	depend := make(map[string]bool)
	for _, r := range actualDepend {
		depend[newNode(r).String()] = true
	}

	increment := -1
	index := make(map[string]int)
	nodeStmt := func(node *pkgTreeNode, sb io.StringWriter) {
		i, exist := index[node.String()]
		if !exist {
			increment++
			i = increment
			index[node.String()] = i
			if depend[node.String()] {
				sb.WriteString(fmt.Sprintf("\tid%d(%s):::red\n", i, node.String()))
			} else {
				sb.WriteString(fmt.Sprintf("\tid%d(%s)\n", i, node.String()))
			}
		}
	}

	repeat := make(map[string]bool)
	return func(level int, node *pkgTreeNode, sb io.StringWriter) {
		tmp := node
		for tmp.parent != nil {
			nodeStmt(tmp, sb)
			nodeStmt(tmp.parent, sb)
			i := index[tmp.String()]
			pi := index[tmp.parent.String()]
			ss := fmt.Sprintf("\tid%d-->id%d\n", pi, i)
			if !repeat[ss] {
				sb.WriteString(ss)
				repeat[ss] = true
			}
			tmp = tmp.parent
		}
	}
}

func dfs(node *pkgTreeNode, level int, handler func(level int, node *pkgTreeNode) (isGoing bool)) {
	if node == nil {
		return
	}
	if !handler(level, node) {
		return
	}
	for _, v := range node.dep {
		dfs(v, level+1, handler)
	}
}
