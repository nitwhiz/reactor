package ecs

import "slices"

type Buffer[T comparable] struct {
	buf      []T
	growSize int
}

func NewBuffer[T comparable](initialSize int, growSize int) *Buffer[T] {
	return &Buffer[T]{
		buf:      make([]T, 0, initialSize),
		growSize: growSize,
	}
}

func (r *Buffer[T]) Sort(cmp func(a, b T) int) {
	slices.SortFunc(r.buf, cmp)
}

func (r *Buffer[T]) Clear() {
	r.buf = r.buf[:0]
}

func (r *Buffer[T]) Grow(by int) {
	if by == 0 {
		return
	}

	newBuffer := make([]T, len(r.buf), len(r.buf)+max(by, r.growSize))

	copy(newBuffer, r.buf)

	r.buf = newBuffer
}

func (r *Buffer[T]) Add(element T) {
	if len(r.buf) == cap(r.buf) {
		r.Grow(r.growSize)
	}

	r.buf = append(r.buf, element)
}

func (r *Buffer[T]) Put(index int, element T) {
	if index >= cap(r.buf) {
		r.Grow(index + 1 - len(r.buf))
	}

	if index >= len(r.buf) {
		r.buf = r.buf[:(index + 1)]
	}

	r.buf[index] = element
}

func (r *Buffer[T]) Append(elements []T) {
	if cap(r.buf)-len(r.buf) <= len(elements) {
		r.Grow(len(elements))
	}

	r.buf = append(r.buf, elements...)
}

func (r *Buffer[T]) Remove(element T) {
	bufLen := len(r.buf)

	found := false
	lastElement := r.buf[bufLen-1]

	if lastElement != element {
		for i, id := range r.buf {
			if id == element {
				r.buf[i] = lastElement
				found = true
				break
			}
		}
	} else {
		found = true
	}

	if found {
		r.buf = r.buf[:bufLen-1]
	}
}

func (r *Buffer[T]) Elements() []T {
	return r.buf
}

func (r *Buffer[T]) Size() int {
	return len(r.buf)
}

func (r *Buffer[T]) At(index int) T {
	return r.buf[index]
}

func (r *Buffer[T]) Last() T {
	return r.buf[len(r.buf)-1]
}

func (r *Buffer[T]) RemoveIndex(index int) {
	bufLen := len(r.buf)
	r.buf[index] = r.buf[bufLen-1]
	r.buf = r.buf[:bufLen-1]
}
