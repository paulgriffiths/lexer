package lexer

import "fmt"

// Error is an interface for lexer error types.
type Error interface {
	error
	implementsError()
}

// RegexError is returned when the lexer cannot compile the
// regular expressions passed to it at creation time.
type RegexError struct {
	rErr error
}

// newRegexError creates a new RegexError.
func newRegexError(err error) Error {
	return RegexError{err}
}

// Error returns a string representation of a RegexError.
func (e RegexError) Error() string {
	return fmt.Sprintf("couldn't compile regex: %v", e.rErr)
}

func (e RegexError) implementsError() {}

// MatchError is returned when the lexer finds input that it cannot
// match against any of its lexeme patterns.
type MatchError struct {
	// Index is the index in the input where the matching failure
	// occurred.
	Index int
}

func newMatchError(index int) Error {
	return MatchError{index}
}

// Error returns a string representation of a MatchError.
func (e MatchError) Error() string {
	return fmt.Sprintf("couldn't match input at position %d", e.Index)
}

func (e MatchError) implementsError() {}

// InputError is returned when the lexer cannot read from its input.
type InputError struct {
	iErr error
}

func newInputError(err error) Error {
	return InputError{err}
}

// Error returns a string representation of an InputError.
func (e InputError) Error() string {
	return fmt.Sprintf("couldn't get input: %v", e.iErr)
}

func (e InputError) implementsError() {}
