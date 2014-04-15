package filterr

import (
	"errors"
	"testing"
)

var MyError = errors.New("foobar")

func IsMyError(err error) bool {
	return err == MyError
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
