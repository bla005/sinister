package sinister

import "fmt"

func max(a, b int) int {
	if a > b {
		return a
	} else if a < b {
		return b
	}
	return a
}

type node struct {
	data   *Route
	left   *node
	right  *node
	height int
}

func newNode(data *Route) *node {
	return &node{
		data:   data,
		left:   nil,
		right:  nil,
		height: 0,
	}
}
func rightRotate(n *node) *node {
	x := n.left
	T2 := x.right

	x.right = n
	n.left = T2

	n.height = max(height(n.left), height(n.right)) + 1
	x.height = max(height(x.left), height(x.right)) + 1
	return x
}
func leftRotate(n *node) *node {
	y := n.right
	T2 := y.left

	y.left = n
	n.right = T2

	n.height = max(height(n.left), height(n.right)) + 1
	y.height = max(height(y.left), height(y.right)) + 1
	return y
}
func insert(n *node, data *Route) *node {
	if len(data.RawPath) == 0 {
		panic("invalid route")
	}
	if n == nil {
		fmt.Println("was nil", data)
		n = newNode(data)
	} else if (data.RawPath[0] < n.data.RawPath[0]) && (data.RawPath != n.data.RawPath) {
		fmt.Println("left", data)
		n.left = insert(n.left, data)
	} else if (data.RawPath[0] > n.data.RawPath[0]) && (data.RawPath != n.data.RawPath) {
		fmt.Println("right", data)
		n.right = insert(n.right, data)
	}
	n.height = max(height(n.left), height(n.right)) + 1
	balance := getBalance(n)

	if balance > 1 && data.RawPath[0] < n.left.data.RawPath[0] {
		return rightRotate(n)
	}
	if balance < -1 && data.RawPath[0] > n.right.data.RawPath[0] {
		return leftRotate(n)
	}
	if balance > 1 && data.RawPath[0] > n.left.data.RawPath[0] {
		n.left = leftRotate(n.left)
		return rightRotate(n)
	}
	if balance < -1 && data.RawPath[0] > n.right.data.RawPath[0] {
		n.right = rightRotate(n.right)
		return leftRotate(n)
	}
	return n
}

func findHeight(n *node) int {
	if n == nil {
		return -1
	}
	return max(findHeight(n.left), findHeight(n.right)) + 1
}

func findSubtreeHeight(n *node) int {
	if n == nil {
		return -1
	}
	return findHeight(n.left) - findHeight(n.right)
}

func height(n *node) int {
	if n == nil {
		return 0
	}
	return n.height
}

func getBalance(n *node) int {
	if n == nil {
		return 0
	}
	return height(n.left) - height(n.right)
}

func findNode(n *node, target string) *Route {
	if n == nil {
		return nil
	}
	if n.data.RawPath == target {
		return n.data
	} else if target[0] < n.data.RawPath[0] {
		return findNode(n.left, target)
	} else {
		return findNode(n.right, target)
	}
}