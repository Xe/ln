package ln

import "fmt"

func ExampleSimpleError() {
	Info(F{"err": fmt.Errorf("This is an Error!!!")}, "fooey", F{"bar": "foo"})
	// Output: THIS FAILS BECAUSE LAZINESS err="This is an Error!!!"
}
