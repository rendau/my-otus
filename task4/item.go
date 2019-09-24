package task4

type Item struct {
	list  *List
	value interface{}
	prev  *Item
	next  *Item
}

func (i *Item) Value() interface{} {
	return i.value
}

func (i *Item) Next() *Item {
	return i.next
}
func (i *Item) Prev() *Item {
	return i.prev
}
