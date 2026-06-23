package goroute

import (
	"strings"
)

// node represents a singe segment of a URL path in our prefix tree
type node struct {
	pattern  string  //full path (e.g., /users/:id), only set on leaf nodes
	part     string  //segment of the path (eg, ":id")
	children []*node //child brances
	isWild   bool    // True if the segment is a dynamic parameter (starts with ':' , '*')
}

type router struct {
	roots    map[string]*node        //separate trees fro each HTTP Method (GET, POST, etc.)
	handlers map[string]RouteHandler //Maps the exact pattern to the function
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]RouteHandler),
	}
}

// parsePath splits the path into processable segments, ignorign trailing slashes
func parsePath(path string) []string {
	vs := strings.Split(path, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)

			// A wildcard '*' means everything subsequent is one big parameter, so we stop parsing
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler RouteHandler) {
	parts := parsePath(pattern)
	key := method + "-" + pattern

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.insert(pattern, parts, 0, r.roots[method])
	r.handlers[key] = handler
}

func (r *router) insert(pattern string, parts []string, height int, n *node) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)

	if child == nil { // part not icluded into child
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	r.insert(pattern, parts, height+1, child)
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePath(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}
	n := r.search(searchParts, 0, root)

	if n != nil {
		parts := parsePath(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) search(parts []string, height int, n *node) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := r.search(parts, height+1, child)
		if result != nil {
			return result
		}
	}
	return nil
}

// ---Trie node helper

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes

}
