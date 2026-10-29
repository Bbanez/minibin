package utils

import "sync"

type FindQuery[Item any] func(item Item) bool
type MapTransform[T any, O any] func(item T) O

func Find[T any](items []T, query FindQuery[T]) Option[T] {
	for _, item := range items {
		if !query(item) {
			continue
		}
		return Some(item)
	}
	return None[T]()
}

func FindIndex[T any](items []T, query FindQuery[T]) int {
	for i, item := range items {
		if !query(item) {
			continue
		}
		return i
	}
	return -1
}

func Filter[T any](items []T, query FindQuery[T]) []T {
	result := []T{}
	for _, item := range items {
		if !query(item) {
			continue
		}
		result = append(result, item)
	}
	return result
}

func MapArr[T any, O any](items []T, transform MapTransform[T, O]) []O {
	result := make([]O, len(items))
	for i, item := range items {
		result[i] = transform(item)
	}
	return result
}

func CloneArr[T any](items []*T) []T {
	var result []T
	for _, item := range items {
		result = append(result, *item)
	}
	return result
}

func ArrayRemoveItemAtPos[T any](items []T, pos int) []T {
	if pos < 0 || pos >= len(items) {
		return items
	}
	return append(items[:pos], items[pos+1:]...)
}

// ----------------------------------------------------------------------------
// ---- Array
// ----------------------------------------------------------------------------

type Array[T any] struct {
	Items []T

	_lock sync.Mutex
}

func Arr[T any](items ...T) Array[T] {
	return Array[T]{
		Items: items,
	}
}

func (a *Array[T]) Push(item T) {
	a._lock.Lock()
	defer a._lock.Unlock()
	a.Items = append(a.Items, item)
}

func (a *Array[T]) PushMany(items ...T) {
	a._lock.Lock()
	defer a._lock.Unlock()
	a.Items = append(a.Items, items...)
}

func (a *Array[T]) Pop() *T {
	a._lock.Lock()
	defer a._lock.Unlock()
	if len(a.Items) == 0 {
		return nil
	}
	last := a.Items[len(a.Items)-1]
	a.Items = a.Items[:len(a.Items)-1]
	return &last
}

func (a *Array[T]) Clear() {
	a._lock.Lock()
	defer a._lock.Unlock()
	a.Items = []T{}
}
