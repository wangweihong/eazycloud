package errors

import (
	"github.com/wangweihong/eazycloud/pkg/sets"
)

// Aggregate represents an object that contains multiple errors, but does not
// necessarily have singular semantic meaning.
type Aggregate interface {
	error
}

// NewAggregate converts a slice of errors into an Aggregate interface, which
// is itself an implementation of the error interface.  If the slice is empty,
// this returns nil.
// It will check if any of the element of input error list is nil, to avoid
// nil pointer panic when call Error().
func NewAggregate(errList ...error) Aggregate {
	if len(errList) == 0 {
		return nil
	}
	// In case of input error list contains nil
	var errs []error
	for _, e := range errList {
		if e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	withStacks := make([]*withStack, 0, len(errs))
	for _, e := range errs {
		withStacks = append(withStacks, FromError(e))
	}

	return aggregate(withStacks)
}

// This helper implements the error and Errors interfaces.  Keeping it private
// prevents people from making an aggregate of 0 errors, which is not
// an error, but does satisfy the error interface.
type aggregate []*withStack

// Error is part of the error interface.
func (agg aggregate) Error() string {
	if len(agg) == 0 {
		// This should never happen, really.
		return ""
	}
	if len(agg) == 1 {
		return agg[0].Error()
	}
	seenErrors := sets.NewString()
	result := ""

	for _, withsStack := range agg {
		msg := withsStack.Error()
		if seenErrors.Has(msg) {
			continue
		}
		seenErrors.Insert(msg)
		if len(seenErrors) > 1 {
			result += ", "
		}
		result += msg
	}

	if len(seenErrors) == 1 {
		return result
	}
	return "[" + result + "]"
}
