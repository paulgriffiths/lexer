package lexer_test

import (
	"fmt"
	"github.com/paulgriffiths/lexer"
	"os"
	"strings"
)

func Example() {
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
}
