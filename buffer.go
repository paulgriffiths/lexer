package lexer

import (
	"unicode"
)

// indexedBuffer represents a byte buffer with a stored
// "current byte" index. The lexer will work by reading
// the entire contents of an io.Reader into this buffer,
// and then attempt to successively match its contents
// with regular expressions representing lexeme patterns.
// The index will represent how much of the input we have
// successfully translated into tokens.
type indexedBuffer struct {
	buffer []byte
	index  int
}

// endOfInput checks if we've reached the end of the buffer.
func (b *indexedBuffer) endOfInput() bool {
	return b.index >= len(b.buffer)
}

// advance advances the index by n bytes.
func (b *indexedBuffer) advance(n int) {
	b.index += n
}

// next returns a slice of the buffer fron the index through
// to the end of the buffer.
func (b *indexedBuffer) next() []byte {
	return b.buffer[b.index:]
}

// current returns the byte at the current index. This should
// not be called if we're at the end of the buffer.
func (b *indexedBuffer) current() byte {
	return b.buffer[b.index]
}

// substring returns, in string format, a slice of the buffer
// of n bytes starting from (and including) the current index.
func (b *indexedBuffer) substring(n int) string {
	return string(b.buffer[b.index : b.index+n])
}

// skipWhitespace advances the current index past any whitespace
// characters. The newline character is treated as whitespace
// if the provided argument is true.
func (b *indexedBuffer) skipWhitespace(skipNewline bool) {
	for !b.endOfInput() {
		r := b.buffer[b.index]
		if (!skipNewline && r == '\n') || !unicode.IsSpace(rune(r)) {
			break
		}
		b.index++
	}
}
