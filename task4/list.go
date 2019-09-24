package task4

type List struct {
	firstItem *Item
	lastItem  *Item
	len       uint64
}

func (l *List) Len() uint64 {
	return l.len
}

func (l *List) First() *Item {
	return l.firstItem
}

func (l *List) Last() *Item {
	return l.lastItem
}

func (l *List) PushFront(v interface{}) (ni *Item) {
	ni = &Item{
		list:  l,
		value: v,
		prev:  nil,
		next:  l.firstItem,
	}

	if l.firstItem == nil {
		l.lastItem = ni
	} else {
		l.firstItem.prev = ni
	}

	l.firstItem = ni

	l.len++

	return
}

func (l *List) PushBack(v interface{}) (ni *Item) {
	ni = &Item{
		list:  l,
		value: v,
		prev:  l.lastItem,
		next:  nil,
	}

	if l.lastItem == nil {
		l.firstItem = ni
	} else {
		l.lastItem.next = ni
	}

	l.lastItem = ni

	l.len++

	return
}

func (l *List) Remove(item *Item) {
	if item != nil && l.len > 0 && item.list == l {
		if item.next == nil {
			l.lastItem = item.prev
		} else {
			item.next.prev = item.prev
		}
		if item.prev == nil {
			l.firstItem = item.next
		} else {
			item.prev.next = item.next
		}

		l.len--

		item.list = nil
		item.next = nil
		item.prev = nil
	}
}
