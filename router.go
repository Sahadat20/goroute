package goroute

type node struct {
	pattern  string  //full path (e.g., /users/:id), only set on leaf nodes
	part     string  //segment of the path (eg, ":id")
	children []*node //child brances
	isWild   bool    // True if the segment is a dynamic parameter (starts with ':' , '*')

}
