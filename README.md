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
