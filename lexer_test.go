package lexer_test

import (
	"github.com/paulgriffiths/lexer"
	"strings"
	"testing"
)

func TestLexerGood(t *testing.T) {
	type lexerTestCase struct {
		lexemes []string
		input   string
		tokens  lexer.TokenList
	}

	testCases := []lexerTestCase{
		{
			// Words and numbers. Note that it will split "ten40"
			// into two tokens, even though they're not separated
			// by a space. We can suppress this behavior by using
			// the \b empty string to match word boundaries
			// (see TestLexerNoMatch for an example).

			[]string{
				"[[:alpha:]]+",
				"[[:digit:]]+",
			},
			"how 2 fail 435 times with 99 ice creams ten40 dog",
			lexer.TokenList{
				lexer.Token{0, "how", 0},
				lexer.Token{1, "2", 4},
				lexer.Token{0, "fail", 6},
				lexer.Token{1, "435", 11},
				lexer.Token{0, "times", 15},
				lexer.Token{0, "with", 21},
				lexer.Token{1, "99", 26},
				lexer.Token{0, "ice", 29},
				lexer.Token{0, "creams", 33},
				lexer.Token{0, "ten", 40},
				lexer.Token{1, "40", 43},
				lexer.Token{0, "dog", 46},
			},
		},
		{
			// Words and numbers and identifiers consisting of an initial
			// letter followed by any number of letters or digits. Note that
			// in the case of words, "[[:alpha:]]+" will be chosen over
			// "[[:alpha:]][[:alnum:]]+" because it appears first in the
			// list, but "ten40" is no longer separated into two tokens
			// because even though "[[:alpha:]][[:alnum:]]+" appears later
			// in the list, the lexer prefers the longest match.

			[]string{
				"[[:alpha:]]+",
				"[[:digit:]]+",
				"[[:alpha:]][[:alnum:]]+",
			},
			"how 2 fail 435 times with 99 ice creams ten40 dog",
			lexer.TokenList{
				lexer.Token{0, "how", 0},
				lexer.Token{1, "2", 4},
				lexer.Token{0, "fail", 6},
				lexer.Token{1, "435", 11},
				lexer.Token{0, "times", 15},
				lexer.Token{0, "with", 21},
				lexer.Token{1, "99", 26},
				lexer.Token{0, "ice", 29},
				lexer.Token{0, "creams", 33},
				lexer.Token{2, "ten40", 40},
				lexer.Token{0, "dog", 46},
			},
		},
		{
			// Words and numbers and identifiers consisting of an initial
			// letter followed by any number of letters or digits. This time
			// "[[:alpha:]][[:alnum:]]+" comes first in the list, so
			// "[[:alpha:]]+" will never be matched as any ambiguities
			// will be resolved in favor of the earliest match.

			[]string{
				"[[:alpha:]][[:alnum:]]+",
				"[[:alpha:]]+",
				"[[:digit:]]+",
			},
			"how 2 fail 435 times with 99 ice creams ten40 dog",
			lexer.TokenList{
				lexer.Token{0, "how", 0},
				lexer.Token{2, "2", 4},
				lexer.Token{0, "fail", 6},
				lexer.Token{2, "435", 11},
				lexer.Token{0, "times", 15},
				lexer.Token{0, "with", 21},
				lexer.Token{2, "99", 26},
				lexer.Token{0, "ice", 29},
				lexer.Token{0, "creams", 33},
				lexer.Token{0, "ten40", 40},
				lexer.Token{0, "dog", 46},
			},
		},
		{
			// Similarly, this will correctly prefer "==" to "=", even
			// though the former appears later in the list of lexemes.
			// Note that we have to escape the parentheses, since they
			// have special meaning to the regular expression engine.

			[]string{
				"[[:digit:]]+", "=", "==", "\\(", "\\)",
			},
			"(32 == 47) = (512 == 681)",
			lexer.TokenList{
				lexer.Token{3, "(", 0},
				lexer.Token{0, "32", 1},
				lexer.Token{2, "==", 4},
				lexer.Token{0, "47", 7},
				lexer.Token{4, ")", 9},
				lexer.Token{1, "=", 11},
				lexer.Token{3, "(", 13},
				lexer.Token{0, "512", 14},
				lexer.Token{2, "==", 18},
				lexer.Token{0, "681", 21},
				lexer.Token{4, ")", 24},
			},
		},
		{
			// When we're dealing with normal mathemtical expressions,
			// we end up having to escape quite a lot of characters,
			// because many of them have special meaning to the regular
			// expression engine. Note that failing to escape a * or a
			// + could potentially put the regular expression engine
			// into an infinite loop, when failing to escape other
			// characters may just cause the compilation of the
			// regular expression to fail. Obviously if we want
			// the characters to have their special meaning
			// (such as with "[[:digit:]]+") then we don't escape them.

			[]string{
				"[[:digit:]]+", "\\+", "-", "\\*", "/", "\\(", "\\)",
			},
			"(3 + 4) * (5 / -6)",
			lexer.TokenList{
				lexer.Token{5, "(", 0},
				lexer.Token{0, "3", 1},
				lexer.Token{1, "+", 3},
				lexer.Token{0, "4", 5},
				lexer.Token{6, ")", 6},
				lexer.Token{3, "*", 8},
				lexer.Token{5, "(", 10},
				lexer.Token{0, "5", 11},
				lexer.Token{4, "/", 13},
				lexer.Token{2, "-", 15},
				lexer.Token{0, "6", 16},
				lexer.Token{6, ")", 17},
			},
		},
		{
			// By default, the lexer will treat the newline character as
			// whitespace, and ignore it other than as a separator
			// between tokens.

			[]string{
				"to", "be", "or", "not",
			},
			"to be\nor not to be",
			lexer.TokenList{
				lexer.Token{0, "to", 0},
				lexer.Token{1, "be", 3},
				lexer.Token{2, "or", 6},
				lexer.Token{3, "not", 9},
				lexer.Token{0, "to", 13},
				lexer.Token{1, "be", 16},
			},
		},
		{
			// But if we include the newline character in the list of
			// lexemes, the lexer will recognize it and return it as
			// a token.

			[]string{
				"to", "be", "or", "not", "\n",
			},
			"to be\nor not to be",
			lexer.TokenList{
				lexer.Token{0, "to", 0},
				lexer.Token{1, "be", 3},
				lexer.Token{4, "\n", 5},
				lexer.Token{2, "or", 6},
				lexer.Token{3, "not", 9},
				lexer.Token{0, "to", 13},
				lexer.Token{1, "be", 16},
			},
		},
		{
			// We can use square brackets in our regular expressions.

			[]string{
				"[ab]+", "c+",
			},
			"abab ccc baa aaa cc baaaa",
			lexer.TokenList{
				lexer.Token{0, "abab", 0},
				lexer.Token{1, "ccc", 5},
				lexer.Token{0, "baa", 9},
				lexer.Token{0, "aaa", 13},
				lexer.Token{1, "cc", 17},
				lexer.Token{0, "baaaa", 20},
			},
		},
		{
			// We can use parentheses in our regular expressions, too.

			[]string{
				"(fr(og|ag)|toad)+", "(bit)+",
			},
			"frogbittoadbitbittoadfrogbitfragfrogbitbitbit",
			lexer.TokenList{
				lexer.Token{0, "frog", 0},
				lexer.Token{1, "bit", 4},
				lexer.Token{0, "toad", 7},
				lexer.Token{1, "bitbit", 11},
				lexer.Token{0, "toadfrog", 17},
				lexer.Token{1, "bit", 25},
				lexer.Token{0, "fragfrog", 28},
				lexer.Token{1, "bitbitbit", 36},
			},
		},
		{
			// For context-free grammars

			[]string{
				"[a-df-zA-Z][[:alnum:]']*",
				"`[^`]+`",
				"\\|",
				":",
				"\n",
				"\\be\\b",
			},
			"S : A' | `terminal` | e\nA' : `another`\n",
			lexer.TokenList{
				lexer.Token{0, "S", 0},
				lexer.Token{3, ":", 2},
				lexer.Token{0, "A'", 4},
				lexer.Token{2, "|", 7},
				lexer.Token{1, "`terminal`", 9},
				lexer.Token{2, "|", 20},
				lexer.Token{5, "e", 22},
				lexer.Token{4, "\n", 23},
				lexer.Token{0, "A'", 24},
				lexer.Token{3, ":", 27},
				lexer.Token{1, "`another`", 29},
				lexer.Token{4, "\n", 38},
			},
		},
	}

	for n, tc := range testCases {
		l, err := lexer.New(tc.lexemes)
		if err != nil {
			t.Errorf("case %d, couldn't create lexer: %v", n+1, err)
			continue
		}

		tokens, err := l.Lex(strings.NewReader(tc.input))
		if err != nil {
			t.Errorf("case %d, couldn't get tokens: %v", n+1, err)
			continue
		}

		if !tokens.Equals(tc.tokens) {
			t.Errorf("case %d, tokens not equals, got %v, want %v",
				n+1, tokens, tc.tokens)
		}
	}
}

func TestLexerBadRegexp(t *testing.T) {
	testCases := [][]string{
		[]string{
			"[[:digit",
		},
		[]string{
			")",
		},
		[]string{
			"(",
		},
	}

	for n, tc := range testCases {
		if _, err := lexer.New(tc); err == nil {
			t.Errorf("case %d, regex unexpectly compiled", n+1)
		} else if _, ok := err.(lexer.RegexError); !ok {
			t.Errorf("case %d, error of unexpected type", n+1)
		}
	}
}

func TestLexerNoMatch(t *testing.T) {
	testCases := []struct {
		input string
		index int
	}{
		{"?", 0},
		{"a!", 1},
		{"abab%abab", 4},
		{"abc 123 abc123", 8},
	}

	l, err := lexer.New([]string{
		"\\b[[:alpha:]]+\\b",
		"\\b[[:digit:]]+\\b",
	})
	if err != nil {
		t.Errorf("couldn't create lexer: %v", err)
		return
	}

	for n, tc := range testCases {
		if _, err := l.Lex(strings.NewReader(tc.input)); err == nil {
			t.Errorf("case %d, regex unexpectly matched", n+1)
		} else if lerr, ok := err.(lexer.MatchError); !ok {
			t.Errorf("case %d, error of unexpected type", n+1)
		} else if lerr.Index != tc.index {
			t.Errorf("case %d, got %d, want %d", n+1, lerr.Index, tc.index)
		}
	}
}
