/*
 * Manage a limited history without duplicates
 */
package radio

import ()

type HistoryNode struct {
	val  string
	prev *HistoryNode
	next *HistoryNode
}

type HistorySet struct {
	head    *HistoryNode
	tail    *HistoryNode
	size    uint64
	maxSize uint64
	index   map[string]*HistoryNode
}

func NewHistorySet(maxSize uint64) HistorySet {
	return HistorySet{
		maxSize: maxSize,
		index:   map[string]*HistoryNode{}}
}

func (hs *HistorySet) ToList() []string {
	history := []string{}
	node := hs.head
	for node != nil {
		history = append(history, node.val)
		node = node.next
	}
	return history
}

func (hs *HistorySet) removeTail() {
	if hs.size == 0 {
		return
	}
	if hs.size > 1 {
		hs.tail.prev.next = nil
	} else {
		hs.head = nil
	}
	delete(hs.index, hs.tail.val)
	hs.tail = hs.tail.prev
}

func (hs *HistorySet) placeAtHead(node *HistoryNode) {
	node.next = hs.head
	hs.head.prev = node
	node.prev = nil
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
	hs.index[s] = node
	if hs.size == 0 {
		hs.head = node
		hs.tail = node
	} else {
		hs.placeAtHead(node)
	}
	if hs.size < hs.maxSize {
		hs.size++
	} else {
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
