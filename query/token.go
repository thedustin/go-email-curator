package query

import "strings"

var trailingSpaceTokenTypes = map[tokenType]bool{
	TokenFulltext:   true,
	TokenFieldValue: true,
	TokenOr:         true,
	TokenGroupEnd:   true,
}

type Token struct {
	kind  tokenType
	value string
}

func NewToken(kind tokenType, value string) Token {
	return Token{kind, value}
}

func (t *Token) Kind() tokenType {
	return t.kind
}

func (t *Token) Value() string {
	return t.value
}

func (t *Token) queryString() string {
	v := t.value

	if t.kind == TokenFieldValue && strings.Contains(v, " ") {
		v = "(" + t.value + ")"
	}

	if !t.isTrailingSpaceToken() {
		return v
	}

	return v + " "
}

func (t *Token) isTrailingSpaceToken() bool {
	_, ok := trailingSpaceTokenTypes[t.kind]

	return ok
}

type TokenList []Token

func (tl TokenList) String() string {
	strs := make([]string, len(tl))

	for i, t := range tl {
		strs[i] = t.queryString()
	}

	return strings.Join(strs, "")
}

func (tl TokenList) Describe() string {
	strs := make([]string, len(tl))

	for i, t := range tl {
		v := string(t.kind)

		if t.isTrailingSpaceToken() {
			v = v + " "
		}

		strs[i] = v
	}

	return strings.Join(strs, "")
}
