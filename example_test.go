package filterr_test

import (
	"errors"
	"fmt"
	"os"

	. "github.com/natefinch/filterr"
)

var (
	SpecificError = errors.New("This is a very specific error!")
	OtherError    = errors.New("Some other error.")
)

func Example() {
	fmt.Println(SpecificError == AllowedError())
	fmt.Println(OtherError == NotAllowedError())
	fmt.Println(NotAllowedError())
	// output:
	// true
	// false
	// Some other error.
}

// This function says it will only ever a SpecificError, or a generic error.
func AllowedError() (err error) {
	defer Returns(&err, Is(SpecificError))

	// Since the returned error is equal to SpecificError, it'll be returned as-
	// is.
	return SpecificError
}

// This function says it will only ever return errors that match os.IsNotExist,
// that are equal to SpecficError, or a generic error.
func NotAllowedError() (err error) {
	defer Returns(&err, os.IsNotExist, Is(SpecificError))

	// Since OtherError neither returns true for os.IsNotExist, nor is equal
	// to SpecificError, it will be anonymized through the ErrorFilter function.
	return OtherError
}
