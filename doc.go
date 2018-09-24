/*
Package lexer implements a simple general-purpose lexical analyzer.

A lexer is created by passing it a slice of strings containing regular
expressions against which potential lexemes should match. After
analyzing its input, the lexer will return a list of tokens. One of
the fields in the token structure is an index identifying which pattern
in this slice of strings matched that individual token, enabling the
user to identify which kind of token it is.

The regular expressions passed as strings will be compiled by the normal
Go regexp package, and can contain any regular expression that that
package considers valid. One caution is that no regular expression
passed should contained a *named* capturing group (which would never be
helpful in any case, since the result of the regexp match is never seen,
only the generated tokens).

In addition, since the strings will be passed verbatim to the regexp
package, any characters in the pattern which may have special meaning
to the regular expression engine should be escaped. For example, if
you want to match a literal left parenthesis, the pattern should be
"\(", or "\\(" in source code, since the left parenthesis would otherwise
be treated as the start of a capturing group by the regular expression
engine.

Whitespace may be used to separate tokens, but is otherwise ignored by
the lexical analyzer. The newline character is treated as whitespace and
similarly ignored, unless it is included (by itself) as one of the
patterns passed to the lexical analyzer at creation time, in which case
each newline character will be returned as a separate token (unless
another pattern embeds a newline character, such as may occur with
multi-line comments in source code.)
*/
package lexer
