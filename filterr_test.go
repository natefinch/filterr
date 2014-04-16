package filterr

import (
	"errors"
	"testing"
)

var MyError = errors.New("foobar")

func IsMyError(err error) bool {
	return err == MyError
}

func AlwaysMyError(_ error) error {
	return MyError
}

func TestIsMatchesSameError(t *testing.T) {
	match := Is(MyError)
	if !match(MyError) {
		t.Error("Failed to match same error")
	}
}

func TestIsWontMatchDifferentError(t *testing.T) {
	match := Is(MyError)
	if match(errors.New("other")) {
		t.Error("Matched wrong error")
	}
}

func TestReturnsAllowsMatchingError(t *testing.T) {
	err := MyError
	Returns(&err, IsMyError)
	if err != MyError {
		t.Error("Error got replaced")
	}
}

func TestReturnsReplacesNonMatchingError(t *testing.T) {
	orig := errors.New("other")
	err := orig
	Returns(&err, IsMyError)
	if err == orig {
		t.Error("Error didn't get replaced")
	}
}

func TestReturnsReturnsNilErrorAsNil(t *testing.T) {
	var err error
	Returns(&err, IsMyError)
	if err != nil {
		t.Error("Error should have been returned as nil, but was non-nil.")
	}
}

func TestReturnsWontBlowUpWithNilPointer(t *testing.T) {
	Returns(nil, IsMyError)
}

func TestMakeReturnFuncHandlesNilMissFilter(t *testing.T) {
	f := MakeReturnFunc(nil, AlwaysMyError)
	orig := errors.New("foo")
	err := orig

	// IsMyError will not match, so should trigger the miss filter, which is nil.
	f(&err, IsMyError)
	if err != orig {
		t.Error("Matching error was changed but should not have been.")
	}
}

func TestMakeReturnFuncHandlesNilMatchFilter(t *testing.T) {
	f := MakeReturnFunc(AlwaysMyError, nil)

	// err will match IsMyError, so should trigger the nil match filter.
	err := MyError
	f(&err, IsMyError)
	if err != MyError {
		t.Error("Matching error was changed but should not have been.")
	}
}

func TestMakeReturnFuncHandlesMissFilter(t *testing.T) {
	f := MakeReturnFunc(AlwaysMyError, nil)
	err := errors.New("foo")

	// IsMyError will not match, so should trigger the miss filter, which should
	// make it into a MyError instead.
	f(&err, IsMyError)
	if err != MyError {
		t.Error("Non-matching error was not changed but should have been.")
	}
}

func TestMakeReturnFuncHandlesMatchFilter(t *testing.T) {
	f := MakeReturnFunc(nil, AlwaysMyError)
	err := errors.New("foo")

	// Is(err) will match, so should trigger the match filter, which should
	// make it into a MyError instead.
	f(&err, Is(err))
	if err != MyError {
		t.Error("Matching error was not changed but should have been.")
	}
}
