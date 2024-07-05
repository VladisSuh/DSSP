package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Color int

const (
	RED Color = iota
	BLACK
)

type NodeRB struct {
	Key        string
	Value      interface{}
	Color      Color
	LeftChild  *NodeRB
	RightChild *NodeRB
	Parent     *NodeRB
}

type RedBlackTree struct {
	Root *NodeRB
}

func NewRedBlackTree() *RedBlackTree {
	return &RedBlackTree{}
}

func (tree *RedBlackTree) Insert(key string, value interface{}) error {
	tree.insertRB(key, value)
	return nil
}

func (tree *RedBlackTree) Get(key string) (interface{}, error) {
	node := tree.searchRB(tree.Root, key)
	if node == nil {
		return nil, fmt.Errorf("key not found")
	}
	return node.Value, nil
}

func (tree *RedBlackTree) GetRange(minValue, maxValue string) ([]string, error) {
	var result []string
	var getRangeHelper func(node *NodeRB, minValue, maxValue string)
	getRangeHelper = func(node *NodeRB, minValue, maxValue string) {
		if node == nil {
			return
		}
		if node.Key >= minValue {
			getRangeHelper(node.LeftChild, minValue, maxValue)
		}
		if node.Key >= minValue && node.Key <= maxValue {
			result = append(result, node.Key)
		}
		if node.Key <= maxValue {
			getRangeHelper(node.RightChild, minValue, maxValue)
		}
	}
	getRangeHelper(tree.Root, minValue, maxValue)
	return result, nil
}

func (tree *RedBlackTree) Update(key string, value interface{}) error {
	node, err := getNodeRB(tree.Root, key)
	if err != nil {
		return err
	}
	node.Value = value
	return nil
}

func (tree *RedBlackTree) Remove(key string) error {
	tree.deleteRB(key)
	return nil
}

func (tree *RedBlackTree) SaveToFile(filename string) error {
	data, err := json.Marshal(tree)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (tree *RedBlackTree) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, tree)
}

func getNodeRB(root *NodeRB, key string) (*NodeRB, error) {
	if root == nil {
		return nil, fmt.Errorf("key not found")
	}
	if key < root.Key {
		return getNodeRB(root.LeftChild, key)
	} else if key > root.Key {
		return getNodeRB(root.RightChild, key)
	} else {
		return root, nil
	}
}

func (tree *RedBlackTree) insertRB(key string, value interface{}) {
	newNode := &NodeRB{
		Key:        key,
		Value:      value,
		Color:      RED,
		LeftChild:  nil,
		RightChild: nil,
		Parent:     nil,
	}
	if tree.Root == nil {
		tree.Root = newNode
	} else {
		tree.insertNodeRB(tree.Root, newNode)
		tree.fixInsertionRB(newNode)
	}
}

func (tree *RedBlackTree) insertNodeRB(root, newNode *NodeRB) {
	if newNode.Key < root.Key {
		if root.LeftChild == nil {
			root.LeftChild = newNode
			newNode.Parent = root
		} else {
			tree.insertNodeRB(root.LeftChild, newNode)
		}
	} else {
		if root.RightChild == nil {
			root.RightChild = newNode
			newNode.Parent = root
		} else {
			tree.insertNodeRB(root.RightChild, newNode)
		}
	}
}

func (tree *RedBlackTree) rotateLeftRB(node *NodeRB) {
	if node == nil || node.RightChild == nil {
		return
	}

	rightChild := node.RightChild
	node.RightChild = rightChild.LeftChild
	if rightChild.LeftChild != nil {
		rightChild.LeftChild.Parent = node
	}
	rightChild.Parent = node.Parent
	if node.Parent == nil {
		tree.Root = rightChild
	} else if node == node.Parent.LeftChild {
		node.Parent.LeftChild = rightChild
	} else {
		node.Parent.RightChild = rightChild
	}
	rightChild.LeftChild = node
	node.Parent = rightChild
}

func (tree *RedBlackTree) rotateRightRB(node *NodeRB) {
	if node == nil || node.LeftChild == nil {
		return
	}

	leftChild := node.LeftChild
	node.LeftChild = leftChild.RightChild
	if leftChild.RightChild != nil {
		leftChild.RightChild.Parent = node
	}
	leftChild.Parent = node.Parent
	if node.Parent == nil {
		tree.Root = leftChild
	} else if node == node.Parent.RightChild {
		node.Parent.RightChild = leftChild
	} else {
		node.Parent.LeftChild = leftChild
	}
	leftChild.RightChild = node
	node.Parent = leftChild
}

func (tree *RedBlackTree) fixInsertionRB(node *NodeRB) {
	for node != nil && node != tree.Root && node.Parent.Color == RED {
		if node.Parent == node.Parent.Parent.LeftChild {
			uncle := node.Parent.Parent.RightChild
			if uncle != nil && uncle.Color == RED {
				node.Parent.Color = BLACK
				uncle.Color = BLACK
				node.Parent.Parent.Color = RED
				node = node.Parent.Parent
			} else {
				if node == node.Parent.RightChild {
					node = node.Parent
					tree.rotateLeftRB(node)
				}
				node.Parent.Color = BLACK
				node.Parent.Parent.Color = RED
				tree.rotateRightRB(node.Parent.Parent)
			}
		} else {
			uncle := node.Parent.Parent.LeftChild
			if uncle != nil && uncle.Color == RED {
				node.Parent.Color = BLACK
				uncle.Color = BLACK
				node.Parent.Parent.Color = RED
				node = node.Parent.Parent
			} else {
				if node == node.Parent.LeftChild {
					node = node.Parent
					tree.rotateRightRB(node)
				}
				node.Parent.Color = BLACK
				node.Parent.Parent.Color = RED
				tree.rotateLeftRB(node.Parent.Parent)
			}
		}
	}
	tree.Root.Color = BLACK
}

func (tree *RedBlackTree) deleteRB(key string) {
	nodeToDelete := tree.searchRB(tree.Root, key)
	if nodeToDelete == nil {
		return
	}

	var child *NodeRB
	if nodeToDelete.LeftChild == nil || nodeToDelete.RightChild == nil {
		child = nodeToDelete
	} else {
		child = tree.successor(nodeToDelete)
	}

	var replacement *NodeRB
	if child.LeftChild != nil {
		replacement = child.LeftChild
	} else {
		replacement = child.RightChild
	}

	if replacement != nil {
		replacement.Parent = child.Parent
	}

	if child.Parent == nil {
		tree.Root = replacement
	} else if child == child.Parent.LeftChild {
		child.Parent.LeftChild = replacement
	} else {
		child.Parent.RightChild = replacement
	}

	if child != nodeToDelete {
		nodeToDelete.Key = child.Key
		nodeToDelete.Value = child.Value
	}

	if child.Color == BLACK && replacement != nil {
		tree.fixDeletionRB(replacement)
	}
}

func (tree *RedBlackTree) searchRB(node *NodeRB, key string) *NodeRB {
	if node == nil || node.Key == key {
		return node
	}

	if node.Key < key {
		return tree.searchRB(node.RightChild, key)
	}
	return tree.searchRB(node.LeftChild, key)
}

func (tree *RedBlackTree) successor(node *NodeRB) *NodeRB {
	if node.RightChild != nil {
		return tree.minimum(node.RightChild)
	}

	parent := node.Parent
	for parent != nil && node == parent.RightChild {
		node = parent
		parent = parent.Parent
	}
	return parent
}

func (tree *RedBlackTree) minimum(node *NodeRB) *NodeRB {
	for node.LeftChild != nil {
		node = node.LeftChild
	}
	return node
}

func (tree *RedBlackTree) fixDeletionRB(node *NodeRB) {
	for node != nil && node != tree.Root && node.Color == BLACK {
		if node == node.Parent.LeftChild {
			sibling := node.Parent.RightChild
			if sibling.Color == RED {
				sibling.Color = BLACK
				node.Parent.Color = RED
				tree.rotateLeftRB(node.Parent)
				sibling = node.Parent.RightChild
			}
			if sibling.LeftChild.Color == BLACK && sibling.RightChild.Color == BLACK {
				sibling.Color = RED
				node = node.Parent
			} else {
				if sibling.RightChild.Color == BLACK {
					sibling.LeftChild.Color = BLACK
					sibling.Color = RED
					tree.rotateRightRB(sibling)
					sibling = node.Parent.RightChild
				}
				sibling.Color = node.Parent.Color
				node.Parent.Color = BLACK
				sibling.RightChild.Color = BLACK
				tree.rotateLeftRB(node.Parent)
				node = tree.Root
			}
		} else {
			sibling := node.Parent.LeftChild
			if sibling.Color == RED {
				sibling.Color = BLACK
				node.Parent.Color = RED
				tree.rotateRightRB(node.Parent)
				sibling = node.Parent.LeftChild
			}
			if sibling.RightChild.Color == BLACK && sibling.LeftChild.Color == BLACK {
				sibling.Color = RED
				node = node.Parent
			} else {
				if sibling.LeftChild.Color == BLACK {
					sibling.RightChild.Color = BLACK
					sibling.Color = RED
					tree.rotateLeftRB(sibling)
					sibling = node.Parent.LeftChild
				}
				sibling.Color = node.Parent.Color
				node.Parent.Color = BLACK
				sibling.LeftChild.Color = BLACK
				tree.rotateRightRB(node.Parent)
				node = tree.Root
			}
		}
	}
	if node != nil {
		node.Color = BLACK
	}
}

type RedBlackCollection struct {
	Tree *RedBlackTree
}

func NewRedBlackCollection() *RedBlackCollection {
	return &RedBlackCollection{
		Tree: NewRedBlackTree(),
	}
}

func (rb *RedBlackCollection) Insert(key string, value interface{}) error {
	return rb.Tree.Insert(key, value)
}

func (rb *RedBlackCollection) Get(key string) (interface{}, error) {
	return rb.Tree.Get(key)
}

func (rb *RedBlackCollection) GetRange(minValue, maxValue string) ([]string, error) {
	return rb.Tree.GetRange(minValue, maxValue)
}

func (rb *RedBlackCollection) Update(key string, value interface{}) error {
	return rb.Tree.Update(key, value)
}

func (rb *RedBlackCollection) Remove(key string) error {
	return rb.Tree.Remove(key)
}

func (rb *RedBlackCollection) SaveToFile(filename string) error {
	data, err := json.Marshal(rb)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (rb *RedBlackCollection) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, rb)
}
