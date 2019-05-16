package errors

import (
	"fmt"
	"math"
	"os"
	"testing"
)

func TestNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal(r)
		}
	}()
	n := (*Tree)(nil)
	n.Add(fmt.Errorf(""))
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
	et.Add(fmt.Errorf("Multiline error:\nvalue is complex"))
	fmt.Fprintf(os.Stdout, "%s\n", et.Error())

	et.Reset()
	if et.IsError() || len(et.errs) > 0 {
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
}

func ExampleTree() {
	et := New("Check error tree")
	for i := 0; i < 2; i++ {
		et.Add(fmt.Errorf("Error case %d", i))
	}
	fmt.Println(et.Error())

	// Output:
	// Check error tree
	// ├──Error case 0
	// └──Error case 1
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
		et.Add(fmt.Errorf("Parameter `f` is NaN"))
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

	// Output:
	// Check input data
	// ├──Parameter `f` is NaN
	// ├──Parameter `i` is less zero
	// └──Parameter `s` is empty
}
