# go-safeslice

SafeSlice is similar to a slice except it is safe for modification
during range-based traversals.

This is achieved by guaranteeing that there are no changes made to
allocations of slices returned by previous calls to `Get`. That is,
only new calls to `Get` will observe the elements added, removed, or
swapped, to the array before that call to `Get`.

IMPORTANT: See the following examples for safe iteration.

The Go language guarantees that the expression passed to `range` is
evaluated only once. Therefore, the following iteration makes a
single call to `Get` before the iteration begins. As such, the
slice being traversed does not observe the added / removed /
swapped elements.

```go
a := NewArray()
for _, x := range a.Get() {
  a.Append(...)  // SAFE
}

a := NewArray()
for _, x := range a.Get() {
  a.Remove(0)  // SAFE
}
```

The following while loop calls `Get` multiple times, therefore the
length of the array is varying between the loops. While this is
correct, it's extremely inefficient because this is forcing a
"snapshot" of the array on each call to the pair (`Get`, `Remove`), so
don't do this.

```go
a := NewArray()
for len(a.Get()) > 0 {
  a.Remove(0)  // INEFFICIENT; DON'T DO
}
```

The following iteration however is NOT safe because each call to
`Get` observes the removed elements.

```go
a := NewArray()
for i := 0; i < len(a.Get()); i++ {
  a.Remove(i)  // WRONG and UNSAFE
}
```

SafeSlice is not safe for concurrent use but it is safe for
non-concurrent modification during traversal.
