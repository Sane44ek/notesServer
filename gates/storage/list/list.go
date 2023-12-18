package list

import (
	"fmt"
	"notes/gates/storage"
	"reflect"
	"sync"
)

type List struct {
	len       int64
	firstNode *node
	mtx       sync.RWMutex
}

func NewList() *List {
	return &List{len: 0, firstNode: nil}
}

func (l *List) Add(data interface{}) (id int64, e error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	newNode := &node{data: data}

	if l.firstNode == nil {
		l.firstNode = newNode
		l.firstNode.index = 0
		l.len++
		return 0, nil
	}

	nod := l.firstNode
	tp := reflect.TypeOf(nod.data)
	if tp != reflect.TypeOf(data) {
		return -1, storage.ErrMismatchType
	}
	id = 0
	for ; nod.next != nil; nod = nod.next {
		id = nod.next.index
	}
	nod.next = newNode
	nod.next.index = id + 1
	l.len++
	return id + 1, nil
}

func (l *List) Print() {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	if l.firstNode == nil {
		fmt.Println("no data")
		return
	}
	for nod := l.firstNode; nod != nil; nod = nod.next {
		fmt.Println(nod.data)
	}
}

func (l *List) Print_All() {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	if l.firstNode == nil {
		fmt.Println("no data")
		return
	}
	for nod := l.firstNode; nod != nil; nod = nod.next {
		fmt.Println(nod.data, nod.index)
	}
}

// Len возвращает длину списка
func (l *List) Len() (len int64) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	return l.len
}

// RemoveByIndex удаляет элемент из списка по индексу
func (l *List) RemoveByIndex(id int64) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.firstNode == nil {
		fmt.Println("no data")
		return
	}
	if id == l.firstNode.index {
		if l.len > 1 {
			l.firstNode = l.firstNode.next
			l.len--
			// l.refresh_indices()
			return
		}
		if l.len == 1 {

			l.firstNode = nil
			l.len = 0

			l = &List{len: 0, firstNode: nil}
			// l.refresh_indices()
			return
		}
	}

	if id < 0 {
		fmt.Println("give positive index")
		return
	}

	var del_nod *node
	flag := false
	for nod := l.firstNode; nod != nil; nod = nod.next {
		if nod.index == id {
			del_nod = nod
			flag = true
		}
	}
	if !flag {
		return
	}

	var prev_nod *node
	for nod := l.firstNode; nod != nil; nod = nod.next {
		if nod.next == del_nod {
			prev_nod = nod
			break
		}
	}
	prev_nod.next = del_nod.next
	l.len--
	// l.refresh_indices()
}

// RemoveByValue удаляет элемент из списка по значению
func (l *List) RemoveByValue(value interface{}) bool {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.firstNode == nil {
		fmt.Println("no data")
		return false
	}
	var id int64 = 0
	for nod := l.firstNode; nod != nil; nod = nod.next {
		if nod.data == value {
			if id == 0 {
				l.firstNode = l.firstNode.next
				l.refresh_indices()
				return true
			}
			if id+1 < l.len {
				l.find_node(id - 1).next = l.find_node(id + 1)
				l.refresh_indices()
				l.len--
				return true
			}
			l.find_node(id - 1).next = nil
			l.refresh_indices()
			l.len--
			return true
		}
		id++
	}
	fmt.Println("not found")
	return false
}

// RemoveAllByValue удаляет все элементы из списка по значению
func (l *List) RemoveAllByValue(value interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.firstNode == nil {
		fmt.Println("no data")
		return
	}
	for {
		res := l.RemoveByValue(value)
		if !res {
			return
		}
		// нужно узнать проблема ли, что выводится на экран 'not found'
	}
}

// GetByIndex возвращает значение элемента по индексу.
//
// Если элемента с таким индексом нет, то возвращается 0 и false.
func (l *List) GetByIndex(id int64) (value interface{}, ok bool) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	nod := l.find_node(id)
	if nod != nil {
		return nod.data, true
	}
	return 0, false
}

// GetByValue возвращает индекс первого найденного элемента по значению.
//
// Если элемента с таким значением нет, то возвращается 0 и false.
func (l *List) GetByValue(value interface{}) (index int64, ok bool) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	if l.firstNode == nil {
		fmt.Println("no data")
		return 0, false
	}
	var id int64 = 0
	for nod := l.firstNode; nod != nil; nod = nod.next {
		if nod.data == value {
			return id, true
		}
		id++
	}
	fmt.Println("not found")
	return 0, false
}

// GetAllByValue возвращает индексы всех найденных элементов по значению
//
// Если элементов с таким значением нет, то возвращается nil и false.
func (l *List) GetAllByValue(value interface{}) (ids []int64, ok bool) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	if l.firstNode == nil {
		fmt.Println("no data")
		return nil, false
	}
	var id int64 = 0
	for nod := l.firstNode; nod != nil; nod = nod.next {
		if nod.data == value {
			ids = append(ids, id)
		}
		id++
	}
	if len(ids) > 0 {
		return ids, true
	}
	fmt.Println("not found")
	return nil, false
}

// GetAll возвращает все элементы списка
//
// Если список пуст, то возвращается nil и false.
func (l *List) GetAll() (values []interface{}, ok bool) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	if l.firstNode == nil {
		fmt.Println("no data")
		return nil, false
	}
	for nod := l.firstNode; nod != nil; nod = nod.next {
		values = append(values, nod.data)
	}
	if len(values) > 0 {
		return values, true
	}
	fmt.Println("not found")
	return nil, false
}

// Clear очищает список
func (l *List) Clear() {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.len = 0
	l.firstNode = nil
}

func (l *List) refresh_indices() {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	var id int64 = 0
	for nod := l.firstNode; nod != nil; nod = nod.next {
		nod.index = id
		id++
	}
}

func (l *List) find_node(index int64) *node {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	for nod := l.firstNode; nod != nil; nod = nod.next {
		if nod.index == index {
			return nod
		}
	}
	return nil
}