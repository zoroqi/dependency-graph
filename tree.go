package main

import (
	"fmt"
	"strings"
)

type pkgTreeNode struct {
	pkg      string
	version  string
	parent   *pkgTreeNode
	dep      []*pkgTreeNode
	circular bool
}

func (m pkgTreeNode) String() string {
	return m.pkg + "@" + m.version
}

func toString2(root *pkgTreeNode) *strings.Builder {
	sb := &strings.Builder{}
	dfs(root, 0, sb)
	return sb
}

func dfs(node *pkgTreeNode, level int, sb *strings.Builder) {
	if node == nil {
		return
	}
	if node.circular {
		sb.WriteString(fmt.Sprintf("%s-%s:circular\n", levelStr(level), node.String()))
	} else {
		sb.WriteString(fmt.Sprintf("%s-%s\n", levelStr(level), node.String()))
	}
	for _, v := range node.dep {
		dfs(v, level+1, sb)
	}

}

type stackNode2 struct {
	pkg   *pkg
	index int // index in traversal
	node  *pkgTreeNode
}

func buildTree(pkg *pkg) *pkgTreeNode {
	stack := make([]*stackNode2, 0, 10)

	stackMap := make(map[string]bool)

	push := func(l *stackNode2) {
		stack = append(stack, l)
		stackMap[l.pkg.String()] = true
	}

	pop := func() *stackNode2 {
		if len(stack) < 0 {
			return nil
		}
		r := stack[len(stack)-1]
		stackMap[r.pkg.String()] = false
		stack = stack[0 : len(stack)-1]
		return r
	}

	top := func() *stackNode2 {
		if len(stack) == 0 {
			return nil
		}
		return stack[len(stack)-1]
	}

	push(&stackNode2{
		pkg:   pkg,
		index: 0,
		node:  newNode(pkg),
	})

	root := top().node

	push2 := func(tmp *stackNode2) {
		n := newNode(tmp.pkg.dep[tmp.index])
		if stackMap[tmp.pkg.dep[tmp.index].String()] {
			n.circular = true
		} else {
			push(&stackNode2{
				pkg:   tmp.pkg.dep[tmp.index],
				index: 0,
				node:  n,
			})
		}
		n.parent = tmp.node
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
		pkg:      pkg.name,
		version:  pkg.ver,
		parent:   nil,
		dep:      nil,
		circular: false,
	}
}
