/*
Package container provides a functional-style compatible immutable list. Users can only append to a list by creating a new list that points to an existing one.
Therefore, any existing list remains immuted. It is, in other words, a prepend-only list.
*/
package container

type ImmutableList struct {
	next  *ImmutableList
	value interface{}
}

func NewImmutableList(next *ImmutableList, value interface{}) *ImmutableList {
	l := new(ImmutableList)
	l.next = next
	l.value = value

	return l
}

func (head *ImmutableList) Next() (next *ImmutableList) {
	return head.next
}

func (head *ImmutableList) Value() (value interface{}) {
	return head.value
}

func (head *ImmutableList) ReverseList() (r *ImmutableList, size int) {
	if head == nil {
		return nil, 0
	}

	r = &ImmutableList{nil, head.Value()}
	for cursor := head; cursor != nil; size, cursor = size+1, cursor.Next() {
		r = NewImmutableList(r, cursor.Value())
	}

	return r, size
}
