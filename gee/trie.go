package gee

import(
	"strings"
	//"fmt"
)

type node struct{
	pattern string //待匹配路由
	part string	//当前结点值
	children []*node //孩子结点
	isWild bool	//是否模糊匹配
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil	//匹配失败
}

func (n *node) matchChildren(part string) []*node {
	var children []*node
	for _, child := range n.children {
		//fmt.Printf("%q,%q\n", child.part, child.iswild)
		if (child.part == part && !child.isWild) || child.isWild {
			children = append(children, child)
		}
	}
	return children
}

func (n *node) insert(pattern string, parts []string, height int){
	if len(parts) == height{
		//将pattern赋值给路由结束节点
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*"){
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		//fmt.Printf("%q", child.pattern)
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}