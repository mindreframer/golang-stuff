package main

import (
	"github.com/emicklei/hopwatch"
)

type node struct {
	label string
	parent *node
	children []node
}

func main() {
	tree := node{label:"parent", children:[]node{node{label:"child"}}}

	// uses go-spew, see https://github.com/davecgh/go-spew
	hopwatch.Dump(tree).Break()
	hopwatch.Dumpf("kids %#+v",tree.children).Break()
}