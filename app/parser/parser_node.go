package main

import (
	"fmt"
	"strings"
)

// dslNode represents a single dslNode in the Abstract Syntax Tree (AST).
// Each dslNode can be a function call, argument, variable reference, or literal value.
type dslNode struct {
	kind     dslNodeKind // The type of node, determining how it should be evaluated
	data     string      // The actual content/value of the node (function name, string value, etc.)
	children []*dslNode  // Child nodes, used for nested function calls and arguments
	named    bool        // Whether this node represents a named argument (e.g. param=value)
	argName  string      // The name of the argument if this is a named argument
	next     *dslNode    // The next node in the sequence, used for chaining statements
}

func (n *dslNode) String() string {
	typ := "any"
	switch n.kind {
	case dsl.nodes.call:
		typ = "func"
	case dsl.nodes.arg:
		typ = "arg"
	case dsl.nodes.varRef:
		typ = "var"
	case dsl.nodes.str:
		typ = "string"
	case dsl.nodes.float:
		typ = "float"
	case dsl.nodes.integer:
		typ = "int"
	case dsl.nodes.boolean:
		typ = "bool"
	case dsl.nodes.assign:
		typ = "assign"
	case dsl.nodes.terminator:
		typ = "end"
	case dsl.nodes.argRef:
		typ = "arg ref"
	}
	return fmt.Sprintf("Node{Type: %s, Value: %s, Children: %v, Named: %t, ArgName: %s}", typ, n.data, n.children, n.named, n.argName)
}

func (n *dslNode) toTree() string {
	var result strings.Builder
	var buildTree func(node *dslNode, prefix string, isLast bool, isRoot bool)

	buildTree = func(node *dslNode, prefix string, isLast bool, isRoot bool) {
		if !isRoot {
			// Add the current node's prefix
			if prefix != "" {
				result.WriteString(prefix)
				if isLast {
					result.WriteString("└── ")
				} else {
					result.WriteString("├── ")
				}
			}
		}

		if node.named {
			if node.kind == dsl.nodes.call {
				fmt.Fprintf(&result, "\x1b[34m%s\x1b[0m: \x1b[33m%s\x1b[0m\n", node.argName, node.data)
			} else {
				fmt.Fprintf(&result, "\x1b[34m%s\x1b[0m: \x1b[32m%s\x1b[0m\n", node.argName, node.data)
			}
		} else {
			if node.kind == dsl.nodes.call {
				fmt.Fprintf(&result, "\x1b[33m%s\x1b[0m\n", node.data)
			} else {
				fmt.Fprintf(&result, "\x1b[32m%s\x1b[0m\n", node.data)
			}
		}

		// Calculate the new prefix for children
		newPrefix := prefix
		if !isRoot {
			if prefix != "" {
				if isLast {
					newPrefix += "    "
				} else {
					newPrefix += "│   "
				}
			} else {
				newPrefix = "    "
			}
		}

		// Process children
		for i, child := range node.children {
			isLastChild := i == len(node.children)-1
			if isRoot {
				if isLastChild {
					result.WriteString("└── ")
				} else {
					result.WriteString("├── ")
				}
			}
			buildTree(child, newPrefix, isLastChild, false)
		}
	}
	buildTree(n, "", true, true)
	return result.String()
}
