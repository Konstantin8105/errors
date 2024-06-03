package errors

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"runtime/debug"
	"testing"
)

// cpu: Intel(R) Xeon(R) CPU           X5550  @ 2.67GHz
// BenchmarkAdd/nil-8         	321649880	       3.692 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAdd/tree.nil-8    	13806802	       99.74 ns/op	      48 B/op	       1 allocs/op
// BenchmarkAdd/tree-8        	11324592	       102.1 ns/op	      48 B/op	       1 allocs/op
//
// BenchmarkAdd/nil-8         	1000000000	      0.6858 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAdd/tree.nil-8    	15273708	       103.7 ns/op	      48 B/op	       1 allocs/op
// BenchmarkAdd/tree-8        	12396004	       97.07 ns/op	      48 B/op	       1 allocs/op
//
// BenchmarkAdd/nil-8         	1000000000	         0.6672 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAdd/tree.nil-8    	506320782	         2.335 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAdd/tree-8        	529980600	         2.066 ns/op	       0 B/op	       0 allocs/op
func BenchmarkAdd(b *testing.B) {
	b.Run("nil", func(b *testing.B) {
		e := New("test.nil")
		for n := 0; n < b.N; n++ {
			e.Add(nil)
			e.Errs = e.Errs[:0]
		}
	})
	b.Run("tree.nil", func(b *testing.B) {
		e := New("test.tree.nil")
		for n := 0; n < b.N; n++ {
			e.Add((*Tree)(nil))
			e.Errs = e.Errs[:0]
		}
	})
	b.Run("tree", func(b *testing.B) {
		e := New("test.tree.nil")
		tr := New("internal")
		for n := 0; n < b.N; n++ {
			e.Add(tr)
			e.Errs = e.Errs[:0]
		}
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("full", func(t *testing.T) {
		te := fmt.Errorf("Unwrap test error")
		e := New("TestUnwrap")
		e.Add(fmt.Errorf("Some tests"))
		e.Add(te)

		if !errors.Is(e, te) {
			t.Fatal("Not found unwrap")
		}
	})
	t.Run("empty", func(t *testing.T) {
		te := fmt.Errorf("Unwrap test error")
		e := New("TestUnwrap")
		if errors.Is(e, te) {
			t.Fatal("Found unwrap")
		}
	})
}

func TestNil(t *testing.T) {
	var buf bytes.Buffer
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("%v\n%s", r, string(debug.Stack()))
		}
	}()
	n := new(Tree)
	n.Add(fmt.Errorf(""))

	// walk
	Walk(n, func(e error) {
		fmt.Fprintf(&buf, "%T %v\n", e, e)
	})
}

func TestWalkNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("%v\n%s", r, string(debug.Stack()))
		}
	}()
	// walk
	Walk((*Tree)(nil), func(e error) {
		fmt.Fprintf(os.Stdout, "%T %v\n", e, e)
	})
}

func ExamplePrint() {
	var et Tree
	for i := 0; i < 10; i++ {
		et.Add(fmt.Errorf("Error %d", i))
		if i%3 == 0 {
			var ett Tree
			for j := 0; j < i/3+1; j++ {
				ett.Add(fmt.Errorf("Inside error %d", j))
				if j%2 == 0 {
					var ettt Tree
					ettt.Name = "Some deep deep errors"
					for k := 0; k < 1+j/2; k++ {
						ettt.Add(fmt.Errorf("Deep error %d", k))
					}
					ett.Add(ettt)
				}
			}
			et.Add(ett)
		}
	}
	et.Add(nil)
	et.Add((error)(nil))
	et.Add(fmt.Errorf("Multiline error:\nvalue is complex"))

	// print error tree
	fmt.Fprintf(os.Stdout, "%s\n", et.Error())

	// walk
	Walk(&et, func(e error) {
		fmt.Fprintf(os.Stdout, "%T %v\n", e, e)
	})

	// reset
	et.Reset()
	if et.IsError() || 0 < len(et.Errs) {
		fmt.Fprintf(os.Stdout, "Reset is not working\n")
	}

	// Output:
	// +
	// ├──Error 0
	// ├──+
	// │  ├──Inside error 0
	// │  └──Some deep deep errors
	// │     └──Deep error 0
	// ├──Error 1
	// ├──Error 2
	// ├──Error 3
	// ├──+
	// │  ├──Inside error 0
	// │  ├──Some deep deep errors
	// │  │  └──Deep error 0
	// │  └──Inside error 1
	// ├──Error 4
	// ├──Error 5
	// ├──Error 6
	// ├──+
	// │  ├──Inside error 0
	// │  ├──Some deep deep errors
	// │  │  └──Deep error 0
	// │  ├──Inside error 1
	// │  ├──Inside error 2
	// │  └──Some deep deep errors
	// │     ├──Deep error 0
	// │     └──Deep error 1
	// ├──Error 7
	// ├──Error 8
	// ├──Error 9
	// ├──+
	// │  ├──Inside error 0
	// │  ├──Some deep deep errors
	// │  │  └──Deep error 0
	// │  ├──Inside error 1
	// │  ├──Inside error 2
	// │  ├──Some deep deep errors
	// │  │  ├──Deep error 0
	// │  │  └──Deep error 1
	// │  └──Inside error 3
	// └──Multiline error:
	//    value is complex
	//
	// *errors.errorString Error 0
	// *errors.errorString Inside error 0
	// *errors.errorString Deep error 0
	// *errors.errorString Error 1
	// *errors.errorString Error 2
	// *errors.errorString Error 3
	// *errors.errorString Inside error 0
	// *errors.errorString Deep error 0
	// *errors.errorString Inside error 1
	// *errors.errorString Error 4
	// *errors.errorString Error 5
	// *errors.errorString Error 6
	// *errors.errorString Inside error 0
	// *errors.errorString Deep error 0
	// *errors.errorString Inside error 1
	// *errors.errorString Inside error 2
	// *errors.errorString Deep error 0
	// *errors.errorString Deep error 1
	// *errors.errorString Error 7
	// *errors.errorString Error 8
	// *errors.errorString Error 9
	// *errors.errorString Inside error 0
	// *errors.errorString Deep error 0
	// *errors.errorString Inside error 1
	// *errors.errorString Inside error 2
	// *errors.errorString Deep error 0
	// *errors.errorString Deep error 1
	// *errors.errorString Inside error 3
	// *errors.errorString Multiline error:
	// value is complex
}

func ExampleTree() {
	et := New("Check error tree")
	for i := 0; i < 2; i++ {
		et.Add(fmt.Errorf("Error case %d", i))
	}
	fmt.Println(et.Error())

	// walk
	Walk(et, func(e error) {
		fmt.Fprintf(os.Stdout, "%T %v\n", e, e)
	})

	// Output:
	// Check error tree
	// ├──Error case 0
	// └──Error case 1
	//
	// *errors.errorString Error case 0
	// *errors.errorString Error case 1
}

type ErrorValue struct {
	ValueName string
	Reason    error
}

func (e ErrorValue) Error() string {
	return fmt.Sprintf("Value `%s`: %v", e.ValueName, e.Reason)
}

func Example() {
	// some input data
	f := math.NaN()
	i := -32
	var s string

	// checking
	var et Tree
	et.Name = "Check input data"
	if math.IsNaN(f) {
		et.Add(ErrorValue{
			ValueName: "f",
			Reason:    fmt.Errorf("is NaN"),
		})
	}
	if f < 0 {
		et.Add(fmt.Errorf("Parameter `f` is negative"))
	}
	if i < 0 {
		et.Add(fmt.Errorf("Parameter `i` is less zero"))
	}
	if s == "" {
		et.Add(fmt.Errorf("Parameter `s` is empty"))
	}

	if et.IsError() {
		fmt.Println(et.Error())
	}

	// walk
	Walk(&et, func(e error) {
		fmt.Fprintf(os.Stdout, "%-25s %v\n", fmt.Sprintf("%T", e), e)
	})

	// Output:
	// Check input data
	// ├──Value `f`: is NaN
	// ├──Parameter `i` is less zero
	// └──Parameter `s` is empty
	//
	// errors.ErrorValue         Value `f`: is NaN
	// *errors.errorString       Parameter `i` is less zero
	// *errors.errorString       Parameter `s` is empty
}
