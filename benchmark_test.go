package filterr

import (
	"errors"
	"os"
	"testing"
)

func BenchmarkDefer(b *testing.B) {
	err := errors.New("text")
	for n := 0; n < b.N; n++ {
		err = DeferReturns(err)
	}
}

func BenchmarkNoDefer(b *testing.B) {
	err := errors.New("text")
	for n := 0; n < b.N; n++ {
		err = NoDeferReturns(err)
	}
}

func DeferReturns(in error) (err error) {
	defer Returns(&err, os.IsNotExist)
	return in
}

func NoDeferReturns(in error) (err error) {
	Returns(&in, os.IsNotExist)
	return in
}
