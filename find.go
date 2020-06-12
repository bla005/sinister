package sinister

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	} else {
// 		return b
// 	}
// 	return a
// }
//
// type node struct {
// 	data  *route
// 	left  *node
// 	right *node
// }
//
// func formatPath(target string) rune {
// 	return strings.ToLower(target)[1]
// }
//
// func isLesserThan(a, b string) bool {
// 	return formatPath(a) <= formatPath(b)
// }
//
// func newNode(data *route) *node {
// 	return &node{
// 		data:  data,
// 		left:  nil,
// 		right: nil,
// 	}
// }
// func insert(n *node, data *route) *node {
// 	if n == nil {
// 		n = newNode(data)
// 	} else if (data.Path[0] <= n.data.Path[0]) && (data.Path != n.data.Path) {
// 		n.left = insert(n.left, data)
// 	} else if (data.Path[0] >= n.data.Path[0]) && (data.Path != n.data.Path) {
// 		n.right = insert(n.right, data)
// 	}
// 	return n
// }
// func findHeight(n *node) int {
// 	if n == nil {
// 		return -1
// 	}
// 	return max(findHeight(n.left), findHeight(n.right)) + 1
// }
//
// func rotate() {}
//
// func findNode(n *node, target *route) *route {
// 	if n == nil {
// 		return nil
// 	}
// 	if n.data.Path == target.Path {
// 		return n.data
// 	} else if isLesserThan(target.Path, n.data.Path) {
// 		return findNode(n.left, target)
// 	} else {
// 		return findNode(n.right, target)
// 	}
// }
