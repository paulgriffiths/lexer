package lexer

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
)

// Lexer implements a general-purpose lexical analyzer.
type Lexer struct {
	lexemes     []string
	regexps     *regexp.Regexp
	skipNewline bool
}

// New creates a new lexer from a slice of strings containing regular
// expressions to match lexemes. Later, the Lex function will return
// a list of tokens with an (id, value) pair. The id will be the index
// in this slice of the pattern that was matched to identify that
// lexeme, so the order is significant.
func New(lexemes []string) (*Lexer, Error) {
	skipNewline := true

	// Build up a combined regular expression for all lexemes
	// so that we may identify them in linear time.

	regexpString := ""
	for i, lexeme := range lexemes {

		// We're going to ignore whitespace between tokens,
		// including newline characters, unless the newline
		// character is specified as one of the lexemes.

		if lexeme == "\n" {
			skipNewline = false
		}
		if i != 0 {
			regexpString += "|"
		}

		// Each lexeme pattern will be a named capturing group
		// in the combined regular expression. We will identify
		// which lexeme pattern we have matched by identifying
		// which capturing group was matched. We have to use
		// named capturing groups here, because if any lexeme
		// pattern contains a parenthesized expression then
		// neither the number of subexpressions nor their ordering
		// will match the slice of lexeme patterns provided to
		// the lexer.

		regexpString += fmt.Sprintf("(?P<%d>^%s)", i, lexeme)
	}

	compiledRegex, err := regexp.Compile(regexpString)
	if err != nil {
		return nil, newRegexError(err)
	}
	compiledRegex.Longest()

	lexer := Lexer{lexemes, compiledRegex, skipNewline}
	return &lexer, nil
}

// Lex lexically analyses the input and returns a list of tokens.
func (l *Lexer) Lex(input io.Reader) (TokenList, Error) {
	bytes, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, newInputError(err)
	}

	buffer := indexedBuffer{bytes, 0}

	list := TokenList{}

	for {
		buffer.skipWhitespace(l.skipNewline)
		if buffer.endOfInput() {
			break
		}

		token, err := l.getNextToken(&buffer)
		if err != nil {
			return nil, err
		}
		list = append(list, token)
	}

	return list, nil
}

// getNextToken gets the next token from a buffer.
func (l *Lexer) getNextToken(b *indexedBuffer) (Token, Error) {

	// Check if there was a match.

	result := l.regexps.FindAllSubmatchIndex(b.next(), 1)
	if len(result) == 0 {
		return Token{-1, string(b.current()), b.index},
			newMatchError(b.index)
	}
	matches := result[0]

	// Loop over the number of subexpressions, which may be different
	// from the number of lexeme patterns initially provided to the
	// lexer if any of the lexeme patterns themselves contain
	// parenthesized capturing groups.

	for i := 0; i < l.regexps.NumSubexp(); i++ {
		beg, end := matches[2*(i+1)], matches[2*(i+1)+1]

		if beg == -1 {

			// There was no match for this subexpression.

			continue
		}

		// There was a match for this subexpression, but we need to
		// check if it has a name, by attempting to convert it to a
		// number (all our named capturing groups are named by
		// sequential numbers). If the user inexplicably named
		// any of the parenthesized capturing groups in their
		// lexeme patterns, then we may be out of luck.

		dn, err := strconv.ParseInt(l.regexps.SubexpNames()[i+1], 10, 32)
		if err != nil {

			// We matched a parenthesized subexpression within
			// one of the lexeme patterns initially provided,
			// and not an entire lexeme pattern, so continue to
			// look for a match.

			continue
		}

		// We found a match, so advance the buffer and return
		// a constructed token.

		token := Token{int(dn), b.substring(end - beg), b.index}
		b.advance(end - beg)
		return token, nil
	}

	// If we got here then we matched the expression but
	// failed to identify the match, which shouldn't happen.

	panic("failed to find regex match index")
}
