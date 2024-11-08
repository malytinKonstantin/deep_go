package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	Key   int
	Value int
	Left  *Node
	Right *Node
}

type OrderedMap struct {
	root *Node
	size int
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		root: nil,
		size: 0,
	}
}

func (m *OrderedMap) Insert(key, value int) {
	newNode := &Node{Key: key, Value: value}
	if m.root == nil {
		m.root = newNode
		m.size++
		return
	}

	current := m.root
	for {
		if key == current.Key {
			current.Value = value // обновляем значение, если ключ уже существует
			return
		}
		if key < current.Key {
			if current.Left == nil {
				current.Left = newNode
				m.size++
				return
			}
			current = current.Left
		} else {
			if current.Right == nil {
				current.Right = newNode
				m.size++
				return
			}
			current = current.Right
		}
	}
}

func (m *OrderedMap) Contains(key int) bool {
	current := m.root
	for current != nil {
		if key == current.Key {
			return true
		}
		if key < current.Key {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	return false
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	inorderTraversal(m.root, action)
}

func inorderTraversal(node *Node, action func(int, int)) {
	if node == nil {
		return
	}
	inorderTraversal(node.Left, action)
	action(node.Key, node.Value)
	inorderTraversal(node.Right, action)
}

func (m *OrderedMap) Erase(key int) {
	m.root = deleteNode(m.root, key, &m.size)
}

func deleteNode(root *Node, key int, size *int) *Node {
	if root == nil {
		return nil
	}

	if key < root.Key {
		root.Left = deleteNode(root.Left, key, size)
	} else if key > root.Key {
		root.Right = deleteNode(root.Right, key, size)
	} else {
		// Нашли узел для удаления
		if root.Left == nil {
			*size--
			return root.Right
		} else if root.Right == nil {
			*size--
			return root.Left
		}

		// Узел имеет два потомка
		minNode := findMin(root.Right)
		root.Key = minNode.Key
		root.Value = minNode.Value
		root.Right = deleteNode(root.Right, minNode.Key, size)
	}
	return root
}

func findMin(node *Node) *Node {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
