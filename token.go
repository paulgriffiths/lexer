package lexer

// Token is a lexical token output by the lexical analyzer.
type Token struct {
	// ID is index of the string slice of lexeme patterns used to
	// create the lexer at which the lexeme pattern used to identify
	// this token is located.
	ID int
	// Value is the actual string value of the lexeme found by the
	// lexical analyzer.
	Value string
	// Index is the position of the input at which the lexeme was
	// found.
	Index int
}

// Equals tests if two tokens are equal.
func (t Token) Equals(other Token) bool {
	return t.ID == other.ID &&
		t.Value == other.Value &&
		t.Index == other.Index
}

// Less tests if a token is less than another token.
func (t Token) Less(other Token) bool {
	if t.Value < other.Value {
		return true
	}
	if t.ID < other.ID {
		return true
	}
	return t.Index < other.Index
}
