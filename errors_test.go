package errors

import (
	"fmt"
	"math"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func TestErrorTree(t *testing.T) {
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
	t.Log(et.Error())

	if err := cupaloy.SnapshotMulti("Tree", et.Error()); err != nil {
		t.Fatalf("error: %s", err)
	}

	et.Reset()
	if et.IsError() || len(et.errs) > 0 {
		t.Fatalf("Reset is not working")
	}
}

func ExampleTree() {
	et := New("Check error tree")
	for i := 0; i < 2; i++ {
		et.Add(fmt.Errorf("Error case %d", i))
	}
	fmt.Println(et.Error())

	// Output:
	// Check error tree
	// ├── Error case 0
	// └── Error case 1
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
	// ├── Parameter `f` is NaN
	// ├── Parameter `i` is less zero
	// └── Parameter `s` is empty
}
