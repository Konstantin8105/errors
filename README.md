[![Go Report Card](https://goreportcard.com/badge/github.com/Konstantin8105/errors)](https://goreportcard.com/report/github.com/Konstantin8105/errors)
[![GoDoc](https://godoc.org/github.com/Konstantin8105/errors?status.svg)](https://godoc.org/github.com/Konstantin8105/errors)
![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)

# errors

Create error tree.

### Installation

```cmd
go get -u github.com/Konstantin8105/errors
```

### Example

```go
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
```

```go
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
```
