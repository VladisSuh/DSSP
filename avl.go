package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Node struct {
	Key    string
	Value  interface{}
	Height int
	Left   *Node
	Right  *Node
}

type AVLTree struct {
	Root *Node
}

func NewAVLTree() *AVLTree {
	return &AVLTree{}
}

func (tree *AVLTree) Insert(key string, value interface{}) error {
	var err error
	tree.Root, err = insert(tree.Root, key, value)
	return err
}

func (tree *AVLTree) Get(key string) (interface{}, error) {
	node, err := getNode(tree.Root, key)
	if err != nil {
		return nil, err
	}
	return node.Value, nil
}

func (tree *AVLTree) GetRange(minValue, maxValue string) ([]string, error) {
	var result []string
	var getRangeHelper func(node *Node, minValue, maxValue string)
	getRangeHelper = func(node *Node, minValue, maxValue string) {
		if node == nil {
			return
		}
		if node.Key >= minValue {
			getRangeHelper(node.Left, minValue, maxValue)
		}
		if node.Key >= minValue && node.Key <= maxValue {
			result = append(result, node.Key)
		}
		if node.Key <= maxValue {
			getRangeHelper(node.Right, minValue, maxValue)
		}
	}
	getRangeHelper(tree.Root, minValue, maxValue)
	return result, nil
}

func (tree *AVLTree) Update(key string, value interface{}) error {
	node, err := getNode(tree.Root, key)
	if err != nil {
		return err
	}
	node.Value = value
	return nil
}

func (tree *AVLTree) Remove(key string) error {
	var err error
	tree.Root, err = deleteNode(tree.Root, key)
	return err
}

func (tree *AVLTree) SaveToFile(filename string) error {
	data, err := json.Marshal(tree)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (tree *AVLTree) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, tree)
}

func height(node *Node) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func rightRotate(y *Node) *Node {
	x := y.Left
	T2 := x.Right

	x.Right = y
	y.Left = T2

	y.Height = max(height(y.Left), height(y.Right)) + 1
	x.Height = max(height(x.Left), height(x.Right)) + 1

	return x
}

func leftRotate(x *Node) *Node {
	y := x.Right
	T2 := y.Left

	y.Left = x
	x.Right = T2

	x.Height = max(height(x.Left), height(x.Right)) + 1
	y.Height = max(height(y.Left), height(y.Right)) + 1

	return y
}

func getBalance(node *Node) int {
	if node == nil {
		return 0
	}
	return height(node.Left) - height(node.Right)
}

func insert(node *Node, key string, value interface{}) (*Node, error) {
	if node == nil {
		return &Node{Key: key, Value: value, Height: 1}, nil
	}

	if key < node.Key {
		var err error
		node.Left, err = insert(node.Left, key, value)
		if err != nil {
			return nil, err
		}
	} else if key > node.Key {
		var err error
		node.Right, err = insert(node.Right, key, value)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("element with this key already exists")
	}

	node.Height = 1 + max(height(node.Left), height(node.Right))

	balance := getBalance(node)

	if balance > 1 && key < node.Left.Key {
		return rightRotate(node), nil
	}

	if balance < -1 && key > node.Right.Key {
		return leftRotate(node), nil
	}

	if balance > 1 && key > node.Left.Key {
		node.Left = leftRotate(node.Left)
		return rightRotate(node), nil
	}

	if balance < -1 && key < node.Right.Key {
		node.Right = rightRotate(node.Right)
		return leftRotate(node), nil
	}

	return node, nil
}

func minValueNode(node *Node) *Node {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

func deleteNode(root *Node, key string) (*Node, error) {
	if root == nil {
		return root, errors.New("element not found")
	}

	if key < root.Key {
		var err error
		root.Left, err = deleteNode(root.Left, key)
		if err != nil {
			return nil, err
		}
	} else if key > root.Key {
		var err error
		root.Right, err = deleteNode(root.Right, key)
		if err != nil {
			return nil, err
		}
	} else {
		if root.Left == nil || root.Right == nil {
			var temp *Node
			if root.Left != nil {
				temp = root.Left
			} else {
				temp = root.Right
			}

			if temp == nil {
				temp = root
				root = nil
			} else {
				*root = *temp
			}
		} else {
			temp := minValueNode(root.Right)
			root.Key = temp.Key
			root.Value = temp.Value
			var err error
			root.Right, err = deleteNode(root.Right, temp.Key)
			if err != nil {
				return nil, err
			}
		}
	}

	if root == nil {
		return root, nil
	}

	root.Height = max(height(root.Left), height(root.Right)) + 1

	balance := getBalance(root)

	if balance > 1 && getBalance(root.Left) >= 0 {
		return rightRotate(root), nil
	}

	if balance > 1 && getBalance(root.Left) < 0 {
		root.Left = leftRotate(root.Left)
		return rightRotate(root), nil
	}

	if balance < -1 && getBalance(root.Right) <= 0 {
		return leftRotate(root), nil
	}

	if balance < -1 && getBalance(root.Right) > 0 {
		root.Right = rightRotate(root.Right)
		return leftRotate(root), nil
	}

	return root, nil
}

func getNode(node *Node, key string) (*Node, error) {
	if node == nil {
		return nil, errors.New("element not found")
	}

	if key < node.Key {
		return getNode(node.Left, key)
	} else if key > node.Key {
		return getNode(node.Right, key)
	} else {
		return node, nil
	}
}

type AVLCollection struct {
	Tree *AVLTree
}

func NewAVLCollection() *AVLCollection {
	return &AVLCollection{
		Tree: NewAVLTree(),
	}
}

func (avl *AVLCollection) Insert(key string, value interface{}) error {
	var err error
	avl.Tree.Root, err = insert(avl.Tree.Root, key, value)
	return err
}

func (avl *AVLCollection) Get(key string) (interface{}, error) {
	node, err := getNode(avl.Tree.Root, key)
	if err != nil {
		return nil, err
	}
	return node.Value, nil
}

func (avl *AVLCollection) GetRange(minValue, maxValue string) ([]string, error) {
	var result []string
	var getRangeHelper func(node *Node, minValue, maxValue string)
	getRangeHelper = func(node *Node, minValue, maxValue string) {
		if node == nil {
			return
		}
		if node.Key >= minValue {
			getRangeHelper(node.Left, minValue, maxValue)
		}
		if node.Key >= minValue && node.Key <= maxValue {
			result = append(result, node.Key)
		}
		if node.Key <= maxValue {
			getRangeHelper(node.Right, minValue, maxValue)
		}
	}
	getRangeHelper(avl.Tree.Root, minValue, maxValue)
	return result, nil
}

func (avl *AVLCollection) Update(key string, value interface{}) error {
	node, err := getNode(avl.Tree.Root, key)
	if err != nil {
		return err
	}
	node.Value = value
	return nil
}

func (avl *AVLCollection) Remove(key string) error {
	var err error
	avl.Tree.Root, err = deleteNode(avl.Tree.Root, key)
	return err
}

func (avl *AVLCollection) SaveToFile(filename string) error {
	data, err := json.Marshal(avl)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (avl *AVLCollection) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, avl)
}
