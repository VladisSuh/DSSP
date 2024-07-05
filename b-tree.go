package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const m = 2

type NodeB struct {
	Keys     []string
	Children []*NodeB
	Leaf     bool
	Value    interface{}
}

type BTree struct {
	Root *NodeB
}

func NewNodeB(leaf bool, value interface{}) *NodeB {
	return &NodeB{
		Keys:     make([]string, 0),
		Children: make([]*NodeB, 0),
		Leaf:     leaf,
		Value:    value,
	}
}

func NewBTree() *BTree {
	return &BTree{
		Root: NewNodeB(true, nil),
	}
}

func (t *BTree) Insert(key string, value interface{}) error {
	root := t.Root
	if len(root.Keys) == (2*m - 1) {
		newRoot := NewNodeB(false, nil)
		newRoot.Children = append(newRoot.Children, root)
		t.Root = newRoot
		t.splitChild(newRoot, 0)
		t.insertNonFull(newRoot, key, value)
	} else {
		t.insertNonFull(root, key, value)
	}
	return nil
}

func (t *BTree) splitChild(node *NodeB, i int) {
	child := node.Children[i]
	newChild := NewNodeB(child.Leaf, nil)
	mid := len(child.Keys) / 2
	splitKey := child.Keys[mid]

	node.Children = append(node.Children[:i+1], append([]*NodeB{newChild}, node.Children[i+1:]...)...)
	node.Keys = append(node.Keys[:i], append([]string{splitKey}, node.Keys[i:]...)...)

	newChild.Keys = append(newChild.Keys, child.Keys[mid+1:]...)
	child.Keys = child.Keys[:mid]

	if !child.Leaf {
		newChild.Children = append(newChild.Children, child.Children[mid+1:]...)
		child.Children = child.Children[:mid+1]
	}
}

func (t *BTree) insertNonFull(node *NodeB, key string, value interface{}) {
	i := len(node.Keys) - 1
	if node.Leaf {
		for i >= 0 && key < node.Keys[i] {
			i--
		}
		node.Keys = append(node.Keys[:i+1], append([]string{key}, node.Keys[i+1:]...)...)
		node.Value = value
	} else {
		for i >= 0 && key < node.Keys[i] {
			i--
		}
		i++
		if len(node.Children[i].Keys) == (2*m - 1) {
			t.splitChild(node, i)
			if key > node.Keys[i] {
				i++
			}
		}
		t.insertNonFull(node.Children[i], key, value)
	}
}

func (t *BTree) Search(key string) *NodeB {
	return t.search(t.Root, key)
}

func (t *BTree) search(node *NodeB, key string) *NodeB {
	if node == nil {
		return nil
	}
	i := 0
	for i < len(node.Keys) && key > node.Keys[i] {
		i++
	}
	if i < len(node.Keys) && key == node.Keys[i] {
		return node
	}
	if node.Leaf {
		return nil
	}
	return t.search(node.Children[i], key)
}

func (t *BTree) Remove(key string) error {
	t.Root = t.delete(t.Root, key)
	if len(t.Root.Keys) == 0 && len(t.Root.Children) == 1 {
		t.Root = t.Root.Children[0]
	}
	return nil
}

func (t *BTree) delete(node *NodeB, key string) *NodeB {
	i := 0
	for i < len(node.Keys) && key > node.Keys[i] {
		i++
	}
	if i < len(node.Keys) && key == node.Keys[i] {
		if node.Leaf {
			t.removeFromLeaf(node, i)
		} else {
			t.removeFromNonLeaf(node, i)
		}
	} else {
		if node.Leaf {
			fmt.Println("Key", key, "not found")
			return node
		}
		flag := (i == len(node.Keys))
		if len(node.Children[i].Keys) < m {
			t.fill(node, i)
		}
		if flag && i > len(node.Keys) {
			node.Children[i-1] = t.delete(node.Children[i-1], key)
		} else {
			node.Children[i] = t.delete(node.Children[i], key)
		}
	}
	return node
}

func (t *BTree) removeFromLeaf(node *NodeB, idx int) {
	copy(node.Keys[idx:], node.Keys[idx+1:])
	node.Keys = node.Keys[:len(node.Keys)-1]
}

func (t *BTree) removeFromNonLeaf(node *NodeB, idx int) {
	key := node.Keys[idx]
	if len(node.Children[idx].Keys) >= m {
		pred := t.getPred(node, idx)
		node.Keys[idx] = pred
		t.delete(node.Children[idx], pred)
	} else if len(node.Children[idx+1].Keys) >= m {
		succ := t.getSucc(node, idx)
		node.Keys[idx] = succ
		t.delete(node.Children[idx+1], succ)
	} else {
		t.merge(node, idx)
		t.delete(node.Children[idx], key)
	}
}

func (t *BTree) getPred(node *NodeB, idx int) string {
	cur := node.Children[idx]
	for !cur.Leaf {
		cur = cur.Children[len(cur.Children)-1]
	}
	return cur.Keys[len(cur.Keys)-1]
}

func (t *BTree) getSucc(node *NodeB, idx int) string {
	cur := node.Children[idx+1]
	for !cur.Leaf {
		cur = cur.Children[0]
	}
	return cur.Keys[0]
}

func (t *BTree) fill(node *NodeB, idx int) {
	if idx != 0 && len(node.Children[idx-1].Keys) >= m {
		t.borrowFromPrev(node, idx)
	} else if idx != len(node.Keys) && len(node.Children[idx+1].Keys) >= m {
		t.borrowFromNext(node, idx)
	} else {
		if idx != len(node.Keys) {
			t.merge(node, idx)
		} else {
			t.merge(node, idx-1)
		}
	}
}

func (t *BTree) borrowFromPrev(node *NodeB, idx int) {
	child := node.Children[idx]
	sibling := node.Children[idx-1]

	child.Keys = append([]string{node.Keys[idx-1]}, child.Keys...)

	if !child.Leaf {
		child.Children = append([]*NodeB{sibling.Children[len(sibling.Children)-1]}, child.Children...)
	}
	node.Keys[idx-1] = sibling.Keys[len(sibling.Keys)-1]
	sibling.Keys = sibling.Keys[:len(sibling.Keys)-1]
	if !sibling.Leaf {
		sibling.Children = sibling.Children[:len(sibling.Children)-1]
	}
}

func (t *BTree) borrowFromNext(node *NodeB, idx int) {
	child := node.Children[idx]
	sibling := node.Children[idx+1]

	child.Keys = append(child.Keys, node.Keys[idx])

	if !child.Leaf {
		child.Children = append(child.Children, sibling.Children[0])
	}
	node.Keys[idx] = sibling.Keys[0]
	sibling.Keys = sibling.Keys[1:]
	if !sibling.Leaf {
		sibling.Children = sibling.Children[1:]
	}
}

func (t *BTree) merge(node *NodeB, idx int) {
	child := node.Children[idx]
	sibling := node.Children[idx+1]

	child.Keys = append(child.Keys, node.Keys[idx])
	child.Keys = append(child.Keys, sibling.Keys...)
	if !child.Leaf {
		child.Children = append(child.Children, sibling.Children...)
	}

	node.Keys = append(node.Keys[:idx], node.Keys[idx+1:]...)
	node.Children = append(node.Children[:idx+1], node.Children[idx+2:]...)
}

func (t *BTree) Get(key string) (interface{}, error) {
	node := t.Search(key)
	if node == nil {
		return nil, fmt.Errorf("key not found")
	}
	return node.Value, nil
}

func (t *BTree) GetRange(minValue, maxValue string) ([]string, error) {
	keysInRange := make([]string, 0)
	t.traverseRange(t.Root, minValue, maxValue, &keysInRange)
	return keysInRange, nil
}

func (t *BTree) traverseRange(node *NodeB, minValue, maxValue string, keysInRange *[]string) {
	if node == nil {
		return
	}

	i := 0
	for i < len(node.Keys) && node.Keys[i] < minValue {
		i++
	}

	for ; i < len(node.Keys); i++ {
		if node.Children[i] != nil {
			t.traverseRange(node.Children[i], minValue, maxValue, keysInRange)
		}
		if node.Keys[i] >= minValue && node.Keys[i] <= maxValue {
			*keysInRange = append(*keysInRange, node.Keys[i])
		}
		if node.Keys[i] > maxValue {
			break
		}
	}

	if node.Children[i] != nil {
		t.traverseRange(node.Children[i], minValue, maxValue, keysInRange)
	}
}

func (t *BTree) Update(key string, value interface{}) error {
	node := t.Search(key)
	if node == nil {
		return fmt.Errorf("key not found")
	}
	node.Value = value
	return nil
}

func (t *BTree) SaveToFile(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (t *BTree) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, t)
}

type BTreeCollection struct {
	Tree *BTree
}

func NewBTreeCollection() *BTreeCollection {
	return &BTreeCollection{
		Tree: NewBTree(),
	}
}

func (bc *BTreeCollection) Insert(key string, value interface{}) error {
	return bc.Tree.Insert(key, value)
}

func (bc *BTreeCollection) Get(key string) (interface{}, error) {
	return bc.Tree.Get(key)
}

func (bc *BTreeCollection) GetRange(minValue, maxValue string) ([]string, error) {
	return bc.Tree.GetRange(minValue, maxValue)
}

func (bc *BTreeCollection) Update(key string, value interface{}) error {
	return bc.Tree.Update(key, value)
}

func (bc *BTreeCollection) Remove(key string) error {
	return bc.Tree.Remove(key)
}

func (bc *BTreeCollection) SaveToFile(filename string) error {
	data, err := json.Marshal(bc)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (bc *BTreeCollection) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, bc)
}
