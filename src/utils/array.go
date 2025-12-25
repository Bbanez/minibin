package utils

type FindQueryRef[Item any] func(item *Item) bool
type FindQuery[Item any] func(item Item) bool
type MapTransform[T any, O any] func(item T) O
type MapTransformRef[T any, O any] func(item *T) *O

func FindRef[T any](items []*T, query FindQueryRef[T]) Option[*T] {
	for i := 0; i < len(items); i++ {
		item := items[i]
		if query(item) {
			return Some(item)
		}
	}
	return None[*T]()
}
func Find[T any](items []T, query FindQuery[T]) Option[T] {
	for i := 0; i < len(items); i++ {
		item := items[i]
		if query(item) {
			return Some(item)
		}
	}
	return None[T]()
}

func FilterRef[T any](items []*T, query FindQueryRef[T]) []*T {
	result := []*T{}
	for i := 0; i < len(items); i++ {
		item := items[i]
		if query(item) {
			result = append(result, item)
		}
	}
	return result
}
func Filter[T any](items []T, query FindQuery[T]) []T {
	result := []T{}
	for i := 0; i < len(items); i++ {
		item := items[i]
		if query(item) {
			result = append(result, item)
		}
	}
	return result
}

func MapRef[T any, O any](items []*T, transform MapTransformRef[T, O]) []*O {
	result := []*O{}
	for i := 0; i < len(items); i++ {
		item := items[i]
		result = append(result, transform(item))
	}
	return result
}
func Map[T any, O any](items []T, transform MapTransform[T, O]) []O {
	result := []O{}
	for i := 0; i < len(items); i++ {
		item := items[i]
		outputItem := transform(item)
		result = append(result, outputItem)
	}
	return result
}

func CloneArray[T any](items []*T) []T {
	var result []T
	for _, item := range items {
		result = append(result, *item)
	}
	return result
}

func ContainsStr(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}
