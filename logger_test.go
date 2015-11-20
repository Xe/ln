package ln

import "fmt"

func ExampleSimpleError() {
	Info(F{"err": fmt.Errorf("This is an Error!!!")}, "fooey", F{"bar": "foo"})
	// Output: Just show me the output
}

func ExampleDebug() {
	oldPri := DefaultLogger.Pri
	defer func() { DefaultLogger.Pri = oldPri }()

	// set priority to Debug
	DefaultLogger.Pri = PriDebug
	Debug(F{"err": fmt.Errorf("This is an Error!!!")})
	// Output: Just show me the output
}
