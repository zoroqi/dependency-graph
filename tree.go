package main

import (
	"fmt"
	"strings"
)

type pkgTreeNode struct {
	name     string
	version  string
	parent   *pkgTreeNode
	dep      []*pkgTreeNode
	circular bool
}

func (m pkgTreeNode) String() string {
	return m.name + "@" + m.version
}

func searchPrint(root *pkgTreeNode, searchName string) *strings.Builder {
	sb := &strings.Builder{}
	dfs(root, 0, func(height int, n *pkgTreeNode) {
		if n.name == searchName {
			reverseLine(n, sb)
		}
	})
	return sb
}

func reversePrint(root *pkgTreeNode) *strings.Builder {
	sb := &strings.Builder{}
	dfs(root, 0, func(height int, n *pkgTreeNode) {
		if len(n.dep) == 0 {
			reverseLine(n, sb)
		}
	})
	return sb
}

func reverseLine(n *pkgTreeNode, sb *strings.Builder) {
	p := n
	for p != nil {
		sb.WriteString(fmt.Sprintf("%s", p.String()))
		if p.circular {
			sb.WriteString(fmt.Sprintf(":circular"))
		}
		p = p.parent
		if p != nil {
			sb.WriteString(fmt.Sprintf(" -> "))
		}
	}
	sb.WriteString("\n")
}

func treePrint(root *pkgTreeNode) *strings.Builder {
	sb := &strings.Builder{}
	dfs(root, 0, func(height int, n *pkgTreeNode) {
		if n.circular {
			sb.WriteString(fmt.Sprintf("%s-%s:circular\n", levelStr(height), n.String()))
		} else {
			sb.WriteString(fmt.Sprintf("%s-%s\n", levelStr(height), n.String()))
		}
	})
	return sb
}

func dfs(node *pkgTreeNode, height int, handler func(level int, node *pkgTreeNode)) {
	if node == nil {
		return
	}
	handler(height, node)
	for _, v := range node.dep {
		dfs(v, height+1, handler)
	}
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
	}
}
