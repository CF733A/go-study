package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type Node struct {
	data *Commit
	prev *Node
	next *Node
}

type DoubleLinkedList struct {
	head *Node // начальный элемент в списке
	tail *Node // последний элемент в списке
	curr *Node // текущий элемент меняется при использовании методов next, prev
	len  int   // количество элементов в списке
}

type LinkedLister interface {
	LoadData(path string) error
	Init(c []Commit)
	Len() int
	SetCurrent(n int) error
	Current() *Node
	Next() *Node
	Prev() *Node
	Insert(n int, c Commit) error
	Push(c Commit) error
	Delete(n int) error
	DeleteCurrent() error
	Index() (int, error)
	GetByIndex(n int) (*Node, error)
	Pop() *Node
	Shift() *Node
	SearchUUID(uuID string) *Node
	Search(message string) *Node
	Reverse() *DoubleLinkedList
}

// LoadData loads data from a JSON file at the given path into the list.
func (d *DoubleLinkedList) LoadData(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var commits []Commit
	if err := json.Unmarshal(bytes, &commits); err != nil {
		return err
	}
	QuickSort(commits)
	d.Init(commits)
	return nil
}

func (d *DoubleLinkedList) Init(c []Commit) {
	d.head = nil
	d.tail = nil
	d.curr = nil
	d.len = 0

	for _, commit := range c {
		d.Push(commit)
	}
}

func (d *DoubleLinkedList) Push(c Commit) error {
	return d.Insert(d.len, c)
}

// Len получение длины списка
func (d *DoubleLinkedList) Len() int {
	return d.len
}

func (d *DoubleLinkedList) SetCurrent(n int) error {
	if n < 0 && n >= d.len {
		return errors.New("out of range dll")
	}

	result, err := d.GetByIndex(n)
	if err != nil {
		return err
	}

	d.curr = result
	return nil
}

// Current получение текущего элемента
func (d *DoubleLinkedList) Current() *Node {
	return d.curr
}

// Next получение следующего элемента
func (d *DoubleLinkedList) Next() *Node {
	if d.curr == nil && d.curr.next == nil {
		return nil
	}
	d.curr = d.curr.next
	return d.curr
}

// Prev получение предыдущего элемента
func (d *DoubleLinkedList) Prev() *Node {
	if d.curr == nil && d.curr.prev == nil {
		return nil
	}
	d.curr = d.curr.prev
	return d.curr
}

// Insert вставка элемента после n элемента

// Insert inserts a new node with commit c at position n.
func (d *DoubleLinkedList) Insert(n int, c Commit) error {
	if n < 0 || n > d.len {
		return errors.New("index out of bounds")
	}
	newNode := &Node{data: &c}
	if n == 0 {
		if d.head == nil {
			d.head = newNode
			d.tail = newNode
		} else {
			newNode.next = d.head
			d.head.prev = newNode
			d.head = newNode
		}
	} else if n == d.len {
		d.tail.next = newNode
		newNode.prev = d.tail
		d.tail = newNode
	} else {
		current := d.head
		for i := 0; i < n; i++ {
			current = current.next
		}
		newNode.next = current
		newNode.prev = current.prev
		current.prev.next = newNode
		current.prev = newNode
	}
	d.len++
	return nil
}

// Delete удаление n элемента
func (d *DoubleLinkedList) Delete(n int) error {
	if n < 0 && n >= d.len {
		return errors.New("out of range DL")
	}

	switch {
	case d.len == 1:
		d.head = nil
		d.curr = nil
		d.tail = nil
	case n == 0:
		d.head = d.head.next
		d.head.prev = nil
		if d.curr == d.head.prev {
			d.curr = d.head
		}
	case n == d.len-1:
		d.tail = d.tail.prev
		d.tail.next = nil
		if d.curr == d.tail.next {
			d.curr = d.tail
		}
	default:
		current := d.head
		for i := 0; i < n; i++ {
			current = current.next
		}
		current.prev.next = current.next
		current.next.prev = current.prev
		if d.curr == current {
			d.curr = current.next
		}
	}
	d.len--
	return nil
}

// DeleteCurrent удаление текущего элемента
func (d *DoubleLinkedList) DeleteCurrent() error {
	if d.curr == nil {
		return errors.New("no curr elem")
	}
	index, err := d.Index()
	if err != nil {
		return err
	}

	d.Delete(index)
	return nil
}

// Index получение индекса текущего элемента
func (d *DoubleLinkedList) Index() (int, error) {
	if d.curr == nil {
		return -1, errors.New("no current elem")
	}
	current := d.head
	for i := 0; i < d.len; i++ {
		if current == d.curr {
			return i, nil
		}
		current = current.next
	}
	return -1, errors.New("current elem not found")
}

func (d *DoubleLinkedList) GetByIndex(n int) (*Node, error) {
	if n < 0 && n >= d.len {
		return nil, errors.New("n - out of range")
	}

	current := d.head
	for i := 0; i < n; i++ {
		current = current.next
	}
	return current, nil
}

// Pop Операция Pop
func (d *DoubleLinkedList) Pop() *Node {
	if d.tail == nil {
		return nil
	}

	tail := d.tail
	if d.len == 1 {
		d.curr = nil
		d.head = nil
		d.tail = nil
	} else {
		d.tail = d.tail.prev
		d.tail.next = nil
		if d.curr == tail {
			d.curr = d.tail
		}
	}

	d.len--
	return tail
}

// Shift операция shift
func (d *DoubleLinkedList) Shift() *Node {
	if d.head == nil {
		return nil
	}

	head := d.head
	if d.len == 1 {
		d.curr = nil
		d.head = nil
		d.tail = nil
	} else {
		d.head = d.head.next
		d.head.prev = nil
		if d.curr == head {
			d.curr = d.head
		}
	}
	d.len--
	return head
}

// SearchUUID поиск коммита по uuid
func (d *DoubleLinkedList) SearchUUID(uuID string) *Node {
	current := d.head
	for current != nil {
		if current.data.UUID == uuID {
			return current
		}
		current = current.next
	}
	return nil
}

// Search поиск коммита по message
func (d *DoubleLinkedList) Search(message string) *Node {
	current := d.head
	for current != nil {
		if current.data.Message == message {
			return current
		}
		current = current.next
	}
	return nil
}

// Reverse возвращает перевернутый список
func (d *DoubleLinkedList) Reverse() *DoubleLinkedList {
	reverse := &DoubleLinkedList{}
	current := d.tail

	for current != nil {
		reverse.Push(*current.data)
		current = current.prev
	}
	return reverse
}

type Commit struct {
	Message string `json:"message"`
	UUID    string `json:"uuid"`
	Date    string `json:"date"`
}

func compareCommitsEss(commit1, commit2 Commit) int {
	date1, _ := time.Parse("2006-01-02", commit1.Date)
	date2, _ := time.Parse("2006-01-02", commit2.Date)

	if date1.Before(date2) {
		return -1
	} else if date1.After(date2) {
		return 1
	} else {
		return 0
	}
}

// QuickSort sorts an array of Commits in ascending order by their Date.
func QuickSort(commits []Commit) {
	if len(commits) < 2 {
		return
	}

	left, right := 0, len(commits)-1

	// Pick a pivot.
	pivotIndex := len(commits) / 2

	// Move the pivot to the right.
	commits[pivotIndex], commits[right] = commits[right], commits[pivotIndex]

	// Pile elements smaller than the pivot on the left.
	for i := range commits {
		if compareCommitsEss(commits[i], commits[right]) < 0 {
			commits[i], commits[left] = commits[left], commits[i]
			left++
		}
	}

	// Place the pivot after the last smaller element.
	commits[left], commits[right] = commits[right], commits[left]

	// Go down the rabbit hole.
	QuickSort(commits[:left])
	QuickSort(commits[left+1:])
}

func GenerateData(numCommits int) []Commit {
	var commits []Commit
	gofakeit.Seed(0) // Initialize the random seed

	// Define how many commits you want to generate
	for i := 0; i < numCommits; i++ {
		commit := Commit{
			Message: gofakeit.Sentence(5),                                                                                                                // Generate a random sentence with 5 words
			UUID:    gofakeit.UUID(),                                                                                                                     // Generate a random UUID
			Date:    gofakeit.DateRange(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)).Format("2006-01-02"), // Generate a random date between 2020 and 2022
		}
		commits = append(commits, commit)
	}

	return commits
}

func main() {
	// Создаем двусвязный список
	dll := &DoubleLinkedList{}

	// Генерируем тестовые данные
	commits := GenerateData(5)

	// Инициализируем список
	dll.Init(commits)

	fmt.Printf("List length: %d\n", dll.Len())

	// Проходим по списку вперед
	fmt.Println("Forward traversal:")
	dll.SetCurrent(0)
	for node := dll.Current(); node != nil; node = dll.Next() {
		fmt.Printf("Date: %s, Message: %s\n", node.data.Date, node.data.Message)
	}

	// Проходим по списку назад
	fmt.Println("\nBackward traversal:")
	dll.SetCurrent(dll.Len() - 1)
	for node := dll.Current(); node != nil; node = dll.Prev() {
		fmt.Printf("Date: %s, Message: %s\n", node.data.Date, node.data.Message)
	}

	// Пример поиска
	found := dll.SearchUUID(commits[2].UUID)
	if found != nil {
		fmt.Printf("\nFound commit: %s\n", found.data.Message)
	}

	// Пример реверса
	reversed := dll.Reverse()
	fmt.Printf("\nReversed list length: %d\n", reversed.Len())

	err := dll.LoadData("test.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d commits from file\n", dll.Len())

}
