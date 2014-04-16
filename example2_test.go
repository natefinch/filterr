package filterr_test

import (
	"fmt"
	"os"

	"github.com/natefinch/filterr"
)

var Returnz = filterr.MakeReturnFunc(filterMiss, nil)

func filterMiss(err error) error {
	return fmt.Errorf("filtered: %s", err)
}

func DemoCustom(in error) (err error) {
	defer Returnz(&err, os.IsNotExist)

	return in
}

func Example_customReturnsFunc() {
	fmt.Println(DemoCustom(fmt.Errorf("Some random error.")))
	// output:
	// filtered: Some random error.
}
