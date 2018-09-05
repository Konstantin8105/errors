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
