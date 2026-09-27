package main

import (
	"fmt"
	"math/rand"
	"time"
)

type User struct {
	ID   int
	Name string
	Age  int
}

type Node struct {
	index int
	data  *User
	left  *Node
	right *Node
}

type BinaryTree struct {
	root *Node
}

func (t *BinaryTree) insert(user *User) *BinaryTree {
	if t.root == nil {
		t.root = &Node{index: user.ID, data: user}
	} else {
		t.root.insert(user)
	}
	return t
}

func (n *Node) insert(user *User) {
	if user.ID < n.index {
		if n.left == nil {
			n.left = &Node{index: user.ID, data: user}
		} else {
			n.left.insert(user)
		}
	} else if user.ID > n.index {
		if n.right == nil {
			n.right = &Node{index: user.ID, data: user}
		} else {
			n.right.insert(user)
		}
	}
}

func (t *BinaryTree) search(key int) *User {
	if t.root == nil {
		return nil
	}

	return t.root.search(key)
}

func (n *Node) search(key int) *User {
	if n == nil {
		return nil
	}

	if key == n.index {
		return n.data
	} else if key < n.index {
		return n.left.search(key)
	} else {
		return n.right.search(key)
	}
}

func generateData(n int) *BinaryTree {
	timeSource := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(timeSource)
	bt := &BinaryTree{}
	for i := 0; i < n; i++ {
		val := rng.Intn(100)
		bt.insert(&User{
			ID:   val,
			Name: fmt.Sprintf("User%d", val),
			Age:  rand.Intn(50) + 20,
		})
	}
	return bt
}

func main() {
	bt := generateData(50)
	user := bt.search(30)
	if user != nil {
		fmt.Printf("Найден пользователь: %+v\n", user)
	} else {
		fmt.Println("Пользователь не найден")
	}
}
