package query

import (
	"fmt"
	"strings"
)

// See https://marianogappa.github.io/software/2019/06/05/lets-build-a-sql-parser-in-go/
// See https://www.youtube.com/watch?v=HxaD_trXwRE

var knownFields = map[string]bool{
	"from":       true,
	"has":        true,
	"in":         true,
	"is":         true,
	"larger":     true,
	"older_than": true,
	"smaller":    true,
	"subject":    true,
	"to":         true,
}

type Parser struct {
	source string
	tokens TokenList

	lastTokenType tokenType
	i             int
	subGroup      byte
}

func NewParser() *Parser {
	p := &Parser{}

	return p
}

func (p *Parser) Result() TokenList {
	return p.tokens
}

func (p *Parser) Reset() {
	p.source = ""
	p.tokens = nil

	p.lastTokenType = TokenStart
	p.i = 0
	p.subGroup = 0
}

func (p *Parser) Parse(query string) error {
	p.Reset()
	p.source = query

	for p.i < len(query) {
		nextToken, _ := p.peek()

		if err := p.processToken(nextToken); err != nil {
			return err
		}

		p.pop()
	}

	if p.subGroup > 0 {
		return fmt.Errorf("consumed all tokens but still has open sub groups")
	}

	return nil
}

func (p *Parser) processToken(t string) error {
	switch p.lastTokenType {
	case TokenField:
		if t != ":" {
			p.tokens[len(p.tokens)-1].kind = TokenFulltext
			p.lastTokenType = TokenFulltext

			return p.processToken(t)
		}

		p.lastTokenType = TokenEqual
		p.tokens = append(p.tokens, Token{TokenEqual, t})

		return nil
	case TokenEqual:
		p.lastTokenType = TokenStart

		if t[0] == '(' {
			// todo: Unescape other "(", ")"
			t = t[1 : len(t)-1]
		}

		p.tokens = append(p.tokens, Token{TokenFieldValue, t})

		return nil
	default:
		if _, ok := knownFields[t]; ok {
			p.lastTokenType = TokenField
			p.tokens = append(p.tokens, Token{TokenField, t})

			return nil
		}

		if t == "OR" {
			p.lastTokenType = TokenOr
			p.tokens = append(p.tokens, Token{TokenOr, t})

			return nil
		}

		if t == "AND" {
			p.lastTokenType = TokenGroupStart

			return nil
		}

		if t == "" {
			p.lastTokenType = TokenGroupEnd
			p.i++
			p.subGroup--
			p.tokens = append(p.tokens, Token{TokenGroupEnd, ")"})

			return nil
		}

		if t[0] == '-' {
			p.tokens = append(p.tokens, Token{TokenNegate, "-"})
			p.lastTokenType = TokenNegate
			p.i++

			nextToken, _ := p.peek()

			return p.processToken(nextToken)
		}

		if t[0] == '(' {
			p.tokens = append(p.tokens, Token{TokenGroupStart, t[0:1]})
			p.lastTokenType = TokenGroupStart
			p.i++
			p.subGroup++

			nextToken, _ := p.peek()

			return p.processToken(nextToken)
		}

		p.tokens = append(p.tokens, Token{TokenFulltext, t})

		return nil
	}

	return fmt.Errorf("unknown state step %q for token: %q", p.lastTokenType, t)
}

func (p *Parser) pop() string {
	peeked, i := p.peek()

	p.i += i
	p.whitespace()

	return peeked
}

func (p *Parser) whitespace() {
	for ; p.i < len(p.source) && p.source[p.i] == ' '; p.i++ {
		// do nothing, we skip all whitespace token by incrementing everything in the for-head
	}
}

func (p *Parser) peek() (string, int) {
	if p.i >= len(p.source) {
		return "", 0
	}

	if p.source[p.i] == ':' {
		return ":", 1
	}

	for field := range knownFields {
		to := min(len(p.source), p.i+len(field))
		t := strings.ToLower(p.source[p.i:to])

		if t == field {
			return t, len(t)
		}
	}

	if p.source[p.i] == '(' {
		return p.lookupNext([]byte{')'})
	}

	t, i := p.lookupNext(p.valueBoundaries())

	if i != 0 { // remove the space from the value
		return t[0:(i - 1)], i - 1
	}

	return t, len(t)
}

func (p *Parser) valueBoundaries() []byte {
	if p.subGroup > 0 {
		return []byte{' ', ')'}
	}

	return []byte{' '}
}

// lookupNext searches for the next occurence of b byte and returns all content (including the b byte), and the length of the content.
// The length will be zero if the b byte was not found.
func (p *Parser) lookupNext(bs []byte) (string, int) {
	i := p.i

	for ; i < len(p.source); i++ {
		for _, b := range bs {
			if p.source[i] == b {
				t := p.source[p.i:(i + 1)]

				return t, len(t)
			}
		}
	}

	return p.source[p.i:i], 0
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
