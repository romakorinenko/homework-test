package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	First  *ListItem
	Last   *ListItem
	Length int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.Length
}

func (l *list) Front() *ListItem {
	return l.First
}

func (l *list) Back() *ListItem {
	return l.Last
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Prev:  nil,
		Next:  nil,
	}

	if l.First == nil {
		l.First = item
		l.Last = item
		l.Length++

		return item
	}

	item.Next = l.First
	l.First = item
	l.First.Next.Prev = item
	l.Length++

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.Last == nil {
		return l.PushFront(v)
	}

	prevLastItem := l.Last

	lastNode := &ListItem{
		Value: v,
		Prev:  l.Last,
		Next:  nil,
	}

	if prevLastItem.Next == nil {
		l.Last = lastNode
	}

	prevLastItem.Next = lastNode
	l.Length++

	return lastNode
}

func (l *list) Remove(i *ListItem) {
	if l.Length == 0 || i == nil {
		return
	}

	if i.Prev == nil {
		l.First = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.Last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.Length--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.Length <= 1 {
		return
	}

	if i != l.First {
		l.PushFront(i.Value)
		l.Remove(i)
	}
}
