package safeslice_test

import (
	"fmt"
	"testing"

	"github.com/jabolopes/go-safeslice"
	"golang.org/x/exp/slices"
)

type myValue struct{}

func ExampleSafeSlice() {
	a := safeslice.New[int]()
	a.Append(1)
	a.Append(2)
	a.Append(3)

	fmt.Printf("%v", a.Get())

	for range a.Get() {
		a.Remove(0)
	}

	// Output: [1 2 3]
}

func TestAppend(t *testing.T) {
	value1 := &myValue{}
	value2 := &myValue{}

	a := safeslice.New[*myValue]()
	a.Append(value1)
	a.Append(value2)

	want := []*myValue{value1, value2}
	if got := a.Get(); !slices.Equal(want, got) {
		t.Errorf("TestAdd: want %v; got %v", want, got)
	}
}

func TestRemove(t *testing.T) {
	value1 := &myValue{}

	a := safeslice.New[*myValue]()
	a.Append(value1)
	a.Remove(0)

	want := []*myValue{}
	if got := a.Get(); !slices.Equal(want, got) {
		t.Errorf("TestRemove: want %v; got %v", want, got)
	}
}

func TestRemoveAppendOnlyAlloc(t *testing.T) {
	value1 := &myValue{}
	value2 := &myValue{}

	a := safeslice.New[*myValue]()
	a.Append(value1)
	a.Append(value2)
	_ = a.Get()
	a.Remove(0)

	want := []*myValue{value2}
	if got := a.Get(); !slices.Equal(want, got) {
		t.Errorf("TestRemove: want %v; got %v", want, got)
	}
}

func TestRemoveAppendOnlyAlloc2(t *testing.T) {
	value1 := &myValue{}
	value2 := &myValue{}

	a := safeslice.New[*myValue]()
	a.Append(value1)
	a.Append(value2)
	_ = a.Get()
	a.Remove(1)

	want := []*myValue{value1}
	if got := a.Get(); !slices.Equal(want, got) {
		t.Errorf("TestRemove: want %v; got %v", want, got)
	}
}

func TestRemoveUntilEmpty(t *testing.T) {
	value1 := &myValue{}
	value2 := &myValue{}

	a := safeslice.New[*myValue]()
	a.Append(value1)
	a.Append(value2)

	a.Remove(1)
	a.Remove(0)

	want := []*myValue{}
	if got := a.Get(); !slices.Equal(want, got) {
		t.Errorf("TestRemove: want %v; got %v", want, got)
	}
}

func TestSwap(t *testing.T) {
	value1 := &myValue{}
	value2 := &myValue{}
	value3 := &myValue{}

	a := safeslice.New[*myValue]()
	a.Append(value1)
	a.Append(value2)
	a.Append(value3)
	a.Swap(1, 2)

	want := []*myValue{value1, value3, value2}
	if got := a.Get(); !slices.Equal(want, got) {
		t.Errorf("TestSwap: want %v; got %v", want, got)
	}
}

func TestSwapAppendOnlyAlloc(t *testing.T) {
	value1 := &myValue{}
	value2 := &myValue{}
	value3 := &myValue{}

	a := safeslice.New[*myValue]()
	a.Append(value1)
	a.Append(value2)
	a.Append(value3)
	_ = a.Get()
	a.Swap(1, 2)

	want := []*myValue{value1, value3, value2}
	if got := a.Get(); !slices.Equal(want, got) {
		t.Errorf("TestSwap: want %v; got %v", want, got)
	}
}

func TestSwapAppendOnlyAlloc2(t *testing.T) {
	value1 := &myValue{}
	value2 := &myValue{}
	value3 := &myValue{}

	a := safeslice.New[*myValue]()
	a.Append(value1)
	a.Append(value2)
	a.Append(value3)
	_ = a.Get()
	a.Swap(0, 2)

	want := []*myValue{value3, value2, value1}
	if got := a.Get(); !slices.Equal(want, got) {
		t.Errorf("TestSwap: want %v; got %v", want, got)
	}
}

func TestSwapSelf(t *testing.T) {
	value1 := &myValue{}

	a := safeslice.New[*myValue]()
	a.Append(value1)
	a.Swap(0, 0)

	want := []*myValue{value1}
	if got := a.Get(); !slices.Equal(want, got) {
		t.Errorf("TestSwap: want %v; got %v", want, got)
	}
}

func TestAppendWhileRange(t *testing.T) {
	value1 := &myValue{}
	value2 := &myValue{}

	a := safeslice.New[*myValue]()
	a.Append(value1)
	a.Append(value2)

	for _, value := range a.Get() {
		a.Append(value)
	}

	want := []*myValue{value1, value2, value1, value2}
	if got := a.Get(); !slices.Equal(want, got) {
		t.Errorf("TestAppendWhileRange: want %v; got %v", want, got)
	}
}

func TestRemoveWhileRange(t *testing.T) {
	value1 := &myValue{}
	value2 := &myValue{}

	a := safeslice.New[*myValue]()
	a.Append(value1)
	a.Append(value2)

	for range a.Get() {
		a.Remove(0)
	}

	want := []*myValue{}
	if got := a.Get(); !slices.Equal(want, got) {
		t.Errorf("TestRemoveWhileRange: want %v; got %v", want, got)
	}
}
