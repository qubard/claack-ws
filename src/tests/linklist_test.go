package tests

import (
	"github.com/qubard/claack-go/lib/ds"
	"testing"
)

func TestLLCreate(t *testing.T) {
	list := ds.CreateLinkedList()

	if list == nil {
		t.Error("Linked list is nil")
	}
}

func TestLLPush(t *testing.T) {
	list := ds.CreateLinkedList()

	var i uint64 = 0
	for i < 100 {
		if list.Len() != i {
			t.Errorf("Invalid list length %d", list.Len())
			return
		}
		list.Push("a" + string(i))
		i++
	}
}

func TestLLRemove(t *testing.T) {
	list := ds.CreateLinkedList()

	var i uint64 = 0
	for i < 5 {
		list.Push("a" + string(i))
		list.Remove("a" + string(i))

		if list.Len() != 0 {
			t.Error("Failed to remove identical node")
			return
		}

		// Implicitly test Contains()
		if list.Contains("a" + string(i)) {
			t.Error("Node still present after removal")
			return
		}
		i++
	}

	i = 0
	for i < 100 {
		list.Push("a" + string(i))
		i++
	}

	i = 0
	for i < 100 {
		if list.Len() != 100-i {
			t.Error("Invalid length after removal")
			return
		}
		list.Remove("a" + string(i))

		if list.Contains("a" + string(i)) {
			t.Error("List contains key after removal")
			return
		}
		i++
	}

	// Test for any potential allocation errors here
	list.Remove("a0")
	list.Remove("c")
}

func TestLLPop(t *testing.T) {
	list := ds.CreateLinkedList()

	list.Push("a")

	p := list.Pop().Value
	if p != "a" {
		t.Error("Invalid pop expecting a")
	}

	if list.Len() != 0 {
		t.Error("Invalid length after pop")
	}

	list.Push("b")
	list.Push("a")

	if list.Len() != 2 {
		t.Error("Invalid length after push")
	}

	p = list.Pop().Value
	if p != "b" {
		t.Error("Invalid second pop expecting b")
	}

	if list.Len() != 1 {
		t.Error("Invalid length after pop")
	}

	p = list.Pop().Value
	if p != "a" {
		t.Error("Invalid third pop expecting a")
	}

	if list.Len() != 0 {
		t.Error("List not 0 after pop")
	}

	// Test for allocation errors
	// Pop on empty list should not crash
	list.Pop()

	if list.Len() != 0 {
		t.Error("List length not 0 after popping while empty")
	}
}
