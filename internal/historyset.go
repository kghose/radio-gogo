/*
 * Manage a limited history without duplicates
 */
package radio

import (
	"iter"
)

type HistoryNode struct {
	val  string
	prev *HistoryNode
	next *HistoryNode
}

type HistorySet struct {
	head    *HistoryNode
	tail    *HistoryNode
	maxSize int
	index   map[string]*HistoryNode
}

func NewHistorySet(maxSize int) HistorySet {
	return HistorySet{
		maxSize: maxSize,
		index:   map[string]*HistoryNode{}}
}

func (hs *HistorySet) removeTail() {
	if hs.tail == nil {
		return
	}
	if hs.tail.prev == nil {
		hs.head = nil
	} else {
		hs.tail.prev.next = nil
	}
	delete(hs.index, hs.tail.val)
	hs.tail = hs.tail.prev
}

func (hs *HistorySet) placeAtHead(node *HistoryNode) {
	if hs.head == nil {
		hs.head = node
		hs.tail = node
		return
	}

	node.prev = nil
	node.next = hs.head
	hs.head.prev = node
	hs.head = node
}

func (hs *HistorySet) detachNode(node *HistoryNode) {
	a := node.prev
	b := node.next
	if a != nil {
		a.next = b
	}
	if b != nil {
		b.prev = a
	}
}

func (hs *HistorySet) moveToHead(node *HistoryNode) {
	hs.detachNode(node)
	hs.placeAtHead(node)
}

func (hs *HistorySet) addToHead(s string) {
	node := &HistoryNode{
		val: s,
	}
	hs.placeAtHead(node)
	hs.index[s] = node
	if len(hs.index) > hs.maxSize {
		hs.removeTail()
	}
}

func (hs *HistorySet) Add(s string) {
	var node *HistoryNode
	var ok bool
	if node, ok = hs.index[s]; ok {
		hs.moveToHead(node)
	} else {
		hs.addToHead(s)
	}
}

func (hs *HistorySet) Len() int {
	return len(hs.index)
}

func (hs *HistorySet) Latest() string {
	if hs.head != nil {
		return hs.head.val
	} else {
		return ""
	}
}

func (hs *HistorySet) ByLatest() iter.Seq[string] {
	return func(yield func(string) bool) {
		node := hs.head
		for node != nil {
			if !yield(node.val) {
				return
			}
			node = node.next
		}
	}
}

func (hs *HistorySet) ByOldest() iter.Seq[string] {
	return func(yield func(string) bool) {
		node := hs.tail
		for node != nil {
			if !yield(node.val) {
				return
			}
			node = node.prev
		}
	}
}
