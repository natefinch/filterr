package filterr

import "errors"

// ErrorFilter is the function that Returns will use to filter out unmatched errors.
// By default, it simply calls errors.New with message of the original error.
var ErrorFilter = func(err *error) {
	*err = errors.New((*err).Error())
}

// Returns checks if err matches any of the given matchers, and if it does not,
// it runs err through ErrFilter.
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
