package safeslice

import "golang.org/x/exp/slices"

// copyDeleteFromArray functionally deletes index i from a. Returns
// the newly allocated array without the removed element.
func copyDeleteFromArray[S ~[]T, T any](s S, i int) S {
	b := make([]T, len(s)-1)
	copy(b[:i], s[:i])

	if i < len(s)-1 {
		copy(b[i:], s[i+1:])
	}

	return b
}

// SafeSlice is similar to a slice except it is safe for modification
// during range-based traversals.
//
// This is done by guaranteeing that there are no changes made to
// allocations of slices returned by previous calls to 'Get'. That is,
// only new calls to 'Get' will observe the elements added, removed,
// or swapped, to the array before that call to 'Get'.
//
// IMPORTANT: See the following examples for safe iteration.
//
// The Go language guarantees that the expression passed to 'range' is
// evaluated only once. Therefore, the following iteration makes a
// single call to 'Get' before the iteration begins. As such, the
// slice being traversed does not observe the added / removed /
// swapped elements.
//
// a := NewSafeSlice()
// for _, x := range a.Get() {
//   a.Append(...)  // SAFE
// }
//
// a := NewSafeSlice()
// for _, x := range a.Get() {
//   a.Remove(0)  // SAFE
// }
//
// The following while loop calls 'Get' multiple times, therefore the
// length of the array is varying between the loops. While this is
// correct, it's extremely inefficient because this is forcing a
// 'snapshot' of the array on each call to the pair (Get, Remove), so
// don't do this.
//
// a := NewSafeSlice()
// for len(a.Get()) > 0 {
//   a.Remove(0)  // INEFFICIENT; DON'T DO
// }
//
// The following iteration however is NOT safe because each call to
// 'Get' observes the removed elements.
//
// a := NewSafeSlice()
// for i := 0; i < len(a.Get()); i++ {
//   a.Remove(i)  // WRONG and UNSAFE
// }
//
// This is not safe for concurrent use. This is safe for
// non-concurrent modification during traversal.
//
// This performs the worst when the calls to 'Get' and 'Remove' are
// interleaved.
type SafeSlice[T any] struct {
	data            []T
	appendOnlyAlloc bool
}

// Append adds a new element to the end of this. It is safe to call
// this during traversal.
func (s *SafeSlice[T]) Append(elem T) {
	s.data = append(s.data, elem)
}

// Remove removes the element at the given index from this. It is safe
// to call this during traversal.
func (s *SafeSlice[T]) Remove(index int) {
	if s.appendOnlyAlloc {
		s.data = copyDeleteFromArray(s.data, index)
		s.appendOnlyAlloc = false
		return
	}

	s.data = slices.Delete(s.data, index, index+1)
}

// Swap swaps the elements at the given indices.
func (s *SafeSlice[T]) Swap(i, j int) {
	if i == j {
		return
	}

	if s.appendOnlyAlloc {
		s.data = append([]T{}, s.data...)
		s.appendOnlyAlloc = false
	}

	swap := s.data[i]
	s.data[i] = s.data[j]
	s.data[j] = swap
}

// Get returns a "snapshot" of the underlying data as a slice. This
// slice can be iterated and this SafeSlice can be modified during
// that iteration.
func (s *SafeSlice[T]) Get() []T {
	s.appendOnlyAlloc = true
	return s.data
}

func New[T any]() *SafeSlice[T] {
	return &SafeSlice[T]{
		nil,   /* data */
		false, /* appendOnlyAlloc */
	}
}
