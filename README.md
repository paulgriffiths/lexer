# Lexer Package Documentation

    package lexer

    import "github.com/paulgriffiths/lexer"

Package lexer implements a simple general-purpose lexical analyzer.

A lexer is created by passing it a slice of strings containing regular
expressions against which potential lexemes should match. After
analyzing its input, the lexer will return a list of tokens. One of the
fields in the token structure is an index identifying which pattern in
this slice of strings matched that individual token, enabling the user
to identify which kind of token it is.

The regular expressions passed as strings will be compiled by the normal
Go regexp package, and can contain any regular expression that that
package considers valid. One caution is that no regular expression
passed should contained a *named* capturing group (which would never be
helpful in any case, since the result of the regexp match is never seen,
only the generated tokens).

In addition, since the strings will be passed verbatim to the regexp
package, any characters in the pattern which may have special meaning to
the regular expression engine should be escaped. For example, if you
want to match a literal left parenthesis, the pattern should be `\(`, or
`"\\("` in source code, since the left parenthesis would otherwise be
treated as the start of a capturing group by the regular expression
engine.

Whitespace may be used to separate tokens, but is otherwise ignored by
the lexical analyzer. The newline character is treated as whitespace and
similarly ignored, unless it is included (by itself) as one of the
patterns passed to the lexical analyzer at creation time, in which case
each newline character will be returned as a separate token (unless
another pattern embeds a newline character, such as may occur with
multi-line comments in source code.)

## Example

```go
names := []string{"Word", "Number", "Punctuation"}
patterns := []string{"[[:alpha:]]+", "[[:digit:]]+", "[\\.,]"}
input := strings.NewReader("20 cats, catch 100 rats.")

lex, err := lexer.New(patterns)
if err != nil {
    fmt.Fprintf(os.Stderr, "couldn't create lexer: %v", err)
    os.Exit(1)
}

tokens, err := lex.Lex(input)
if err != nil {
    fmt.Fprintf(os.Stderr, "couldn't lex input: %v", err)
    os.Exit(1)
}

for _, t := range tokens {
    fmt.Printf("%-11s : %-7q - found at index %d\n",
        names[t.ID], t.Value, t.Index)
}

// Output:
// Number      : "20"    - found at index 0
// Word        : "cats"  - found at index 3
// Punctuation : ","     - found at index 7
// Word        : "catch" - found at index 9
// Number      : "100"   - found at index 15
// Word        : "rats"  - found at index 19
// Punctuation : "."     - found at index 23
```

# Types

```go
type Error interface {
    error
    // contains filtered or unexported methods
}
```

`Error` is an interface for lexer error types.

```go
type InputError struct {
    // contains filtered or unexported fields
}
```

`InputError` is returned when the lexer cannot read from its input.

```go
func (e InputError) Error() string
```


`Error` returns a string representation of an `InputError`.

```go
type Lexer struct {
    // contains filtered or unexported fields
}
```

`Lexer` implements a general-purpose lexical analyzer.

```go
func New(lexemes []string) (*Lexer, Error)
```


`New` creates a new lexer from a slice of strings containing regular
expressions to match lexemes. Later, the `Lex` function will return a list
of tokens with an (id, value) pair. The id will be the index in this
slice of the pattern that was matched to identify that lexeme, so the
order is significant.

```go
func (l *Lexer) Lex(input io.Reader) (TokenList, Error)
```


`Lex` lexically analyses the input and returns a list of tokens.

```go
type MatchError struct {
    // Index is the index in the input where the matching failure
    // occurred.
    Index int
}
```

`MatchError` is returned when the lexer finds input that it cannot match
against any of its lexeme patterns.

```go
func (e MatchError) Error() string
```


`Error` returns a string representation of a `MatchError`.

```go
type RegexError struct {
    // contains filtered or unexported fields
}
```

`RegexError` is returned when the lexer cannot compile the regular
expressions passed to it at creation time.

```go
func (e RegexError) Error() string
```


`Error` returns a string representation of a `RegexError`.

```go
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
```

`Token` is a lexical token output by the lexical analyzer.

```go
func (t Token) Equals(other Token) bool
```


`Equals` tests if two tokens are equal.

```go
func (t Token) Less(other Token) bool
```


`Less` tests if a token is less than another token.

```go
type TokenList []Token
```

`TokenList` is a list of lexical tokens.

```go
func (t TokenList) Equals(other TokenList) bool
```


`Equals` tests if two token lists are equal.

```go
func (t TokenList) IsEmpty() bool
```


`IsEmpty` checks if the list is empty.

```go
func (t TokenList) Len() int
```


`Len` returns the number of tokens in the list.

```go
func (t TokenList) Less(i, j int) bool
```


`Less` returns true if list[i] < list[j].

```go
func (t TokenList) Swap(i, j int)
```


`Swap` swaps tokens i and j in the list.
