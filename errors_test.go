package errors

import (
	"fmt"
	"math"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func TestErrorTree(t *testing.T) {
	var et ErrorTree
	for i := 0; i < 10; i++ {
		et.Add(fmt.Errorf("Error %d", i))
		if i%3 == 0 {
			var ett ErrorTree
			for j := 0; j < i/3+1; j++ {
				ett.Add(fmt.Errorf("Inside error %d", j))
				if j%2 == 0 {
					var ettt ErrorTree
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
}

func ExampleError() {
	var et ErrorTree
	et.Name = "Check error tree"
	for i := 0; i < 2; i++ {
		et.Add(fmt.Errorf("Error case %d", i))
	}
	fmt.Println(et.Error())

	// Output:
	// Check error tree
	// ├── Error case 0
	// └── Error case 1
}

func ExampleIsError() {
	// some input data
	var f float64 = math.NaN()
	var i int = -32
	var s string = ""

	// chacking
	var et ErrorTree
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
