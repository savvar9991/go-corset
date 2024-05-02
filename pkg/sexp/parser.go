package sexp

import (
	"errors"
	"unicode"
)

// Parse a given string into an S-expression, or return an error if the string
// is malformed.
func Parse(s string) (SExp, error) {
	p := &Parser{s}
	// Parse the input
	sExp, err := p.Parse()
	// Sanity check everything was parsed
	if err == nil && p.text != "" {
		return nil, errors.New("unexpected string remainder")
	}

	return sExp, err
}

// Parser represents a parser in the process of parsing a given string into one
// or more S-expressions.
type Parser struct {
	text string
}

// NewParser constructs a new instance of Parser
func NewParser(text string) *Parser {
	return &Parser{
		text: text,
	}
}

// Parse a given string into an S-Expression, or produce an error.
func (p *Parser) Parse() (SExp, error) {
	token := p.Next()

	if token == "" {
		return nil, nil
	} else if token == ")" {
		return nil, errors.New("unexpected end-of-list")
	} else if token == "(" {
		var elements []SExp

		for p.Lookahead(0) != ")" {
			// Parse next element
			element, err := p.Parse()
			if err != nil {
				return nil, err
			}
			// Continue around!
			elements = append(elements, element)
		}
		// Consume right-brace
		p.Next()
		// Done
		return &List{elements}, nil
	}

	return &Symbol{token}, nil
}

// Next extracts the next token from a given string.
func (p *Parser) Next() string {
	if p.text == "" {
		return ""
	}

	switch p.text[0] {
	case '(', ')':
		// List begin / end
		token := p.text[0:1]
		p.text = p.text[1:]

		return token
	case ' ', '\n':
		// Whitespace
		p.text = p.text[1:]
		return p.Next()
	case ';':
		// Comment
		return p.parseComment()
	}
	// Symbol
	return p.parseSymbol()
}

// Lookahead and see what punctuation is next.
func (p *Parser) Lookahead(i int) string {
	if len(p.text) > i {
		switch p.text[i] {
		case '(', ')', ';':
			return p.text[0:1]
		case ' ', '\n':
			return p.Lookahead(i + 1)
		default:
			return ""
		}
	}

	return ""
}

func (p *Parser) parseSymbol() string {
	// Parse token
	i := len(p.text)

	for j, c := range p.text {
		if c == ')' || unicode.IsSpace(c) {
			i = j
			break
		}
	}
	// Reached end of token
	token := p.text[0:i]
	p.text = p.text[i:]

	return token
}

func (p *Parser) parseComment() string {
	// Parse token
	i := len(p.text)

	for j, c := range p.text {
		if c == '\n' {
			i = j
			break
		}
	}
	// Skipped comment
	p.text = p.text[i:]
	// Look for next token
	return p.Next()
}
