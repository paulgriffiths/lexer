package lexer

// TokenList is a list of lexical tokens.
type TokenList []Token

// Len returns the number of tokens in the list.
func (t TokenList) Len() int {
	return len(t)
}

// IsEmpty checks if the list is empty.
func (t TokenList) IsEmpty() bool {
	return len(t) == 0
}

// Equals tests if two token lists are equal.
func (t TokenList) Equals(other TokenList) bool {
	if len(t) != len(other) {
		return false
	}
	for n := range t {
		if t[n] != other[n] {
			return false
		}
	}
	return true
}

// Less returns true if list[i] < list[j].
func (t TokenList) Less(i, j int) bool {
	if t[i].Value < t[j].Value {
		return true
	}
	return t[i].ID < t[j].ID
}

// Swap swaps tokens i and j in the list.
func (t TokenList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
