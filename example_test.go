package filterr_test

import (
	"fmt"

	"github.com/natefinch/filterr"
)

var (
	SpecificError = fmt.Errorf("This is a very specific error!")
	OtherError    = fmt.Errorf("Some other error.")

	Returns, Is = filterr.Returns, filterr.Is
)

func Demo(in error) (err error) {
	defer Returns(&err, Is(SpecificError))

	return in
}

func Example() {
	fmt.Println(SpecificError == Demo(SpecificError))
	fmt.Println(OtherError == Demo(OtherError))
	fmt.Println(Demo(OtherError))
	// output:
	// true
	// false
	// Some other error.
}
