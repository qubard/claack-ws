package ds

type LinkNode struct {
	Next  *LinkNode
	Prev  *LinkNode
	Value string
}

type LinkedList struct {
	len   uint64
	dummy *LinkNode
	tail  *LinkNode
	m     map[string]*LinkNode // Assume unique keys
}

func (list *LinkedList) Push(value string) {
	tmp := LinkNode{}
	tmp.Value = value

	list.tail.Next = &tmp
	tmp.Prev = list.tail

	list.tail = &tmp
	list.m[value] = &tmp
	list.len++
}

func (list *LinkedList) Len() uint64 {
	return list.len
}

func (list *LinkedList) Contains(value string) bool {
	_, present := list.m[value]
	return present
}

func (list *LinkedList) Remove(value string) {
	if curr, present := list.m[value]; present {
		if curr.Next != nil {
			curr.Next.Prev = curr.Prev
		}
		list.len--
		if curr.Next != nil {
			curr.Prev.Next = curr.Next.Next
		}
		delete(list.m, curr.Value)
	}
}

// Pop the node immediately after the dummy node
func (list *LinkedList) Pop() *LinkNode {
	ret := list.dummy.Next

	// No elements
	if ret == nil {
		return ret
	}

	if ret == list.tail { // Removed the tail from the "head"
		list.tail = ret.Prev
	}

	list.dummy.Next = ret.Next
	if ret.Next != nil {
		ret.Next.Prev = list.dummy
	}

	list.len--
	delete(list.m, ret.Value)
	return ret
}

func CreateLinkedList() *LinkedList {
	dummy := &LinkNode{}
	return &LinkedList{
		dummy: dummy,
		tail:  dummy,
		m:     make(map[string]*LinkNode),
	}
}
