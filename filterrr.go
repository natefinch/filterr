// Package filterr makes it easy to filter the errors a function returns.
//
// Purpose
//
// It is idiomatic in Go to simply pass up errors that get returned from
// functions you call, however, this can have unintended consequences.  If your
// package consistently returns a specific type of error, even if it is
// undocumented, consumers of your paxckage may write code that relies on that
// specific error.  If you later change the implementation such that a different
// error is returned, their code will break.
//
// Code Enforced Error Returns
//
// This package solves this problem by pre-emptively declaring exactly what
// errors a function can return, and enforcing that with code.  Any errors which
// the function is not defined as returning will get wrapped in a generic error
// which masks the type of the original error.
//
// And Documentation
//
// Further, the code that enforces the error returns also serves as
// documentation in the code that is easy to find, since it always resides in a
// defer statement right at the top of the function body like this:
//
//   func TryIt() (err error) {
//   	defer Returns(&err, os.IsNotExist, Is(MyError))
//   	...
//   	return err
//   }
//
// The above function is defined as only returning errors which either match
// os.IsNotExist, or errors that == MyError (a custom error type).  Any other
// errors which the function returns will get wrapped in a errors.New error with
// the error message from the wrapped error.
package filterr

import "errors"

// ErrorFilter is the function that Returns will use to filter out unmatched
// errors. By default, it simply calls errors.New with message of the original
// error, but you can assign whatever function you want to it in order to
// customize your error filtering.
var ErrorFilter = func(err *error) {
	*err = errors.New((*err).Error())
}

// Returns checks if err matches any of the given matchers, and if it does not,
// it runs err through ErrorFilter.
func Returns(err *error, matchers ...func(error) bool) {
	if err == nil {
		return
	}

	for _, match := range matchers {
		if match(*err) {
			return
		}
	}
	ErrorFilter(err)
}

// Is is a helper function that generates a function that can be passed into
// Returns, which returns true if the passed-in error is equal to any of the
// errors given to Is.
func Is(errors ...error) func(error) bool {
	return func(err error) bool {
		for _, match := range errors {
			if match == err {
				return true
			}
		}
		return false
	}
}
