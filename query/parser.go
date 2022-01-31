package query

import (
	"fmt"
	"strings"
)

// See https://marianogappa.github.io/software/2019/06/05/lets-build-a-sql-parser-in-go/
// See https://www.youtube.com/watch?v=HxaD_trXwRE

type step string

const (
	stepStart      step = "^"
	stepGroupStart step = "("
	stepGroupEnd   step = ")"
	stepField      step = "%field%"
	stepEqual      step = ":"
	stepFieldValue step = "%value%"
	stepEnd        step = "$"
	stepOr         step = "OR"
)

var allowedFields = map[string]bool{
	"from":       true,
	"older_than": true,
	"subject":    true,
	"to":         true,
}

type token struct {
	kind  step
	value string
}

type Parser struct {
	source string
	result *Query

	step step
	i    int
}

type Query struct{}

func NewParser() *Parser {
	p := &Parser{}

	return p
}

func (p *Parser) Parse(query string) (*Query, error) {
	p.source = query
	p.step = stepStart

	for p.i < len(query) {
		nextToken, _ := p.peek()

		if err := p.processToken(nextToken); err != nil {
			return nil, err
		}

		p.pop()
	}

	return nil, nil
}

func (p *Parser) processToken(t string) error {
	fmt.Println(t)

	switch p.step {
	case stepField:
		if t != ":" {
			return fmt.Errorf("expected colon \":\" after a field but got %q", t)
		}

		p.step = stepEqual
		return nil
	case stepEqual:
		p.step = stepFieldValue

		return nil
	default:
		if _, ok := allowedFields[t]; ok {
			p.step = stepField

			return nil
		}
	}

	return fmt.Errorf("unknown state step %q for token: %q", p.step, t)
}

func (p *Parser) pop() string {
	peeked, i := p.peek()

	p.i += i
	p.whitespace()

	return peeked
}

func (p *Parser) whitespace() {
	for ; p.i < len(p.source) && p.source[p.i] == ' '; p.i++ {
		// do nothing, we increment everything in the for-head
	}
}

func (p *Parser) peek() (string, int) {
	if p.i >= len(p.source) {
		return "", 0
	}

	if p.source[p.i] == ':' {
		return ":", 1
	}

	for field := range allowedFields {
		to := min(len(p.source), p.i+len(field))
		t := strings.ToLower(p.source[p.i:to])

		if t == field {
			return t, len(t)
		}
	}

	if p.source[p.i] == '(' {
		return p.lookupNext(')')
	}

	t, i := p.lookupNext(' ')

	if i != 0 { // remove the space from the value
		return t[0:(i - 1)], i - 1
	}

	return t, len(t)
}

// lookupNext searches for the next occurence of b byte and returns all content (including the b byte), and the length of the content.
// The length will be zero if the b byte was not found.
func (p *Parser) lookupNext(b byte) (string, int) {
	i := p.i

	for ; i < len(p.source); i++ {
		if p.source[i] == b {
			t := p.source[p.i:(i + 1)]

			return t, len(t)
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
