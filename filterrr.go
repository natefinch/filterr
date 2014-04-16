// Package filterr makes it easy to filter the errors a function returns.
//
// Purpose
//
// In Go code, it is idiomatic to simply pass up errors that get returned from
// functions you call (the ubuquitous if err != nil { return err }).  However,
// this practice can have unintended consequences.  If your package consistently
// returns a specific type of error, even if it is undocumented, consumers of
// your package may write code that relies on that specific error.  If you later
// change the implementation such that a different error is returned, even if
// the actual behavior of the code is identical, their code will break.
//
// Enforced Error Specification
//
// This package solves the above problem by declaring exactly what errors a
// function can return, and enforcing that with code.  Any errors which the
// function is not defined as returning will get wrapped in a generic error
// which masks the type of the original error, preventing any callers of your
// function from depending on the error type.
//
// Also Documentation
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
// The above function declares itself as only returning errors which either
// satisfy os.IsNotExist, or are equal to MyError (a custom error type).  Any
// other errors which the function returns will get wrapped in an anonymous error using errors.New.
//
// Customizing
//
// For projects that want to do something other than a simple anonymization of
// errors, you can customize the transform functions that get called with the
// errors your function returns using MakeReturnFunc.
package filterr

import "errors"

// Returns sanitizes the errors returned from a function.  It checks if err
// matches any of the given matchers, if it does, it is returned as-is.  If it
// does not match, it returns a new anonymous error that reuses the original
// error's Error() string.  If err is nil or points to a nil error, Returns is a
// no-op.
func Returns(err *error, matchers ...func(error) bool) {
	if err == nil || *err == nil {
		return
	}

	for _, m := range matchers {
		if m(*err) {
			return
		}
	}
	*err = errors.New((*err).Error())
}

// MakeReturnFunc allows you to use custom transform functions to create your
// own Returns function.  match will be used to transform errors that match
// the matchers provided to Returns, and miss will be used to transform
// errors which do not match the matchers provided for Returns.  A nil function
// in either case will cause the original error to be returned unmodified.
//
// match exists mainly for applications that want to wrap all returned errors
// using a third party error handling package.  If you are not using such a
// package, it should generally be nil.
func MakeReturnFunc(miss, match func(error) error) (Returns func(*error, ...func(error) bool)) {
	return func(err *error, matchers ...func(error) bool) {
		if err == nil || *err == nil {
			return
		}

		for _, m := range matchers {
			if m(*err) {
				if match != nil {
					*err = match(*err)
				}
				return
			}
		}
		if miss != nil {
			*err = miss(*err)
		}
	}
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
