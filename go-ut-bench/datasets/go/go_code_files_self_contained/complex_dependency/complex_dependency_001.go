package main

type TreeNode struct {
	value int
	left  *TreeNode
	right *TreeNode
}

func NewTreeNode(value int) *TreeNode {
	return &TreeNode{value: value, left: nil, right: nil}
}

type BinarySearchTree struct {
	root *TreeNode
}

func NewBinarySearchTree() *BinarySearchTree {
	return &BinarySearchTree{root: nil}
}

func (bst *BinarySearchTree) Insert(value int) {
	bst.root = insertRecursive(bst.root, value)
}

func insertRecursive(node *TreeNode, value int) *TreeNode {
	if node == nil {
		return NewTreeNode(value)
	}
	if value < node.value {
		node.left = insertRecursive(node.left, value)
	} else {
		node.right = insertRecursive(node.right, value)
	}
	return node
}

func (bst *BinarySearchTree) Search(value int) bool {
	return searchRecursive(bst.root, value)
}

func searchRecursive(node *TreeNode, value int) bool {
	if node == nil {
		return false
	}
	if node.value == value {
		return true
	}
	if value < node.value {
		return searchRecursive(node.left, value)
	}
	return searchRecursive(node.right, value)
}

func (bst *BinarySearchTree) InorderTraversal() []int {
	result := []int{}
	inorderRecursive(bst.root, &result)
	return result
}

func inorderRecursive(node *TreeNode, result *[]int) {
	if node != nil {
		inorderRecursive(node.left, result)
		*result = append(*result, node.value)
		inorderRecursive(node.right, result)
	}
}