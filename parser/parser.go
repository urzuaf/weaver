package parser

import (
	"errors"
	"fmt"
	"strconv"
	"text/scanner"
	"weaver/lexer"
)

type Parser struct {
	l *lexer.Lexer

	buf struct {
		tok lexer.Token
		lit string
		pos scanner.Position
		n   int
	}
}

func NewParser(lexer *lexer.Lexer) *Parser {
	return &Parser{l: lexer}
}

// consume the next token from lexer or bufer
func (p *Parser) scan() (tok lexer.Token, lit string, pos scanner.Position) {
	if p.buf.n > 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit, p.buf.pos
	}
	return p.l.Next()
}

// peek at the next token without consuming it
func (p *Parser) peek() lexer.Token {
	if p.buf.n > 0 {
		return p.buf.tok
	}
	p.buf.tok, p.buf.lit, p.buf.pos = p.l.Next()
	p.buf.n = 1
	return p.buf.tok
}

// entry
func (p *Parser) Parse() (*FileNode, error) {
	file := &FileNode{}

	//iterating till EOF
	for p.peek() != lexer.TokenEOF {
		item, err := p.parseItem()
		if err != nil {
			return nil, err
		}
		file.Items = append(file.Items, item)
	}
	return file, nil
}

// main processing function
func (p *Parser) parseItem() (Node, error) {
	tok, lit, pos := p.scan()
	if tok != lexer.TokenIdentifier {
		return nil, fmt.Errorf("error in %s: expected an identifier, but got '%s'", pos, lit)
	}

	switch p.peek() {
	//Check for named or anonymous blocks
	case lexer.TokenString:
		return p.parseBlock(tok, lit, pos)

	case lexer.TokenLBrace:
		return p.parseBlock(tok, lit, pos)

	// Assignment
	case lexer.TokenEqual, lexer.TokenColon:
		return p.parseAssignment(tok, lit, pos)

	default:
		peekTok, peekLit, peekPos := p.buf.tok, p.buf.lit, p.buf.pos
		return nil, fmt.Errorf("error in %s: unexpected token '%s' (type %v) after identifier '%s'", peekPos, peekLit, peekTok, lit)
	}
}

func (p *Parser) parseBlock(typeTok lexer.Token, typeLit string, typePos scanner.Position) (*BlockNode, error) {
	block := &BlockNode{
		Type: typeLit,
		Name: "",
	}

	if p.peek() == lexer.TokenString {
		_, litName, _ := p.scan()
		block.Name = litName
	}

	//consume '{'
	tok, lit, pos := p.scan()
	if tok != lexer.TokenLBrace {
		return nil, fmt.Errorf("error in %s, expected '{' to start block body, got %s", pos, lit)
	}

	body := &FileNode{}

	for p.peek() != lexer.TokenRBrace && p.peek() != lexer.TokenEOF {
		item, err := p.parseItem()
		if err != nil {
			return nil, err
		}
		body.Items = append(body.Items, item)
	}

	block.Body = body

	//consume '}'
	tok, lit, pos = p.scan()
	if tok != lexer.TokenRBrace {
		return nil, fmt.Errorf("error in %s, expected '}' to close block body, got %s", pos, lit)
	}

	return block, nil
}

func (p *Parser) parseAssignment(keyTok lexer.Token, keyLit string, KeyPos scanner.Position) (*AssignmentNode, error) {
	assign := &AssignmentNode{}

	assign.Key = keyLit

	//consume '=' or ':'
	tok, lit, pos := p.scan()
	if tok != lexer.TokenEqual && tok != lexer.TokenColon {
		return nil, fmt.Errorf("error in %s, expected '=' or ':', got %s", pos, lit)
	}

	// parse value
	val, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	assign.Value = val

	return assign, nil
}

func (p *Parser) parseList() (Node, error) {
	// consume '['
	tok, lit, pos := p.scan()
	if tok != lexer.TokenLBracket {
		return nil, fmt.Errorf("error in %s, expected '[' to start list, got %s", pos, lit)
	}

	list := &ListLiteral{}

	for p.peek() != lexer.TokenRBracket && p.peek() != lexer.TokenEOF {
		var item Node
		var err error

		nextToken := p.peek()
		// if we find an identifier, we parse a block
		if nextToken == lexer.TokenIdentifier {
			identTok, identLit, identPos := p.scan()
			item, err = p.parseBlock(identTok, identLit, identPos)
		} else {
			// else, we parse a value
			item, err = p.parseValue()
		}

		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, item)

		if p.peek() == lexer.TokenComma {
			p.scan()
		} else if p.peek() != lexer.TokenRBracket {
			tok, lit, pos := p.scan()
			return nil, fmt.Errorf("error in %s: expected ']' or ',' in list, but got '%s' (Token: %v)", pos, lit, tok)
		}
	}

	// Consume ']'
	tok, lit, pos = p.scan()
	if tok != lexer.TokenRBracket {
		return nil, fmt.Errorf("error in %s, expected ']' to close list, got %s (Token: %v)", pos, lit, tok)
	}

	return list, nil
}

func (p *Parser) parseNumber() (Node, error) {
	tok, lit, pos := p.scan()

	if tok != lexer.TokenNumber {
		return nil, fmt.Errorf("error in %s, expected number, got %s", pos, lit)
	}
	value, err := strconv.ParseFloat(lit, 64)

	if err != nil {
		return nil, fmt.Errorf("error in %s, invalid number format: %s", pos, lit)
	}

	// Check for magnitude unit
	if p.peek() == lexer.TokenIdentifier {

		peekPosition := p.buf.pos
		if peekPosition.Line == pos.Line {
			_, unitLit, _ := p.scan()
			return &MagnitudeNode{Value: value, Unit: unitLit}, nil
		}

	}

	return &NumberLiteral{Value: value}, nil

}

func (p *Parser) parseReference() (Node, error) {
	var ref []string

	tok, lit, pos := p.scan()
	if tok != lexer.TokenIdentifier {
		// unreachable, but for safety
		return nil, fmt.Errorf("error in %s: expected an identifier, but got '%s'", pos, lit)
	}

	ref = append(ref, lit)
	initialLine := pos.Line // save first reference

	//search for a dot
	for p.peek() == lexer.TokenDot {
		// 3. Check that the dot is on the same line
		if p.buf.pos.Line != initialLine {
			break // The dot is on another line
		}

		p.scan() // Consume .

		// after the dot, we expect another identifier
		tok, lit, pos = p.scan()
		if tok != lexer.TokenIdentifier {
			return nil, fmt.Errorf("error in %s: expected an identifier after '.', but got '%s'", pos, lit)
		}

		// 5. Check that this new identifier is also on the same line
		if pos.Line != initialLine {
			return nil, fmt.Errorf("error in %s: reference cannot span multiple lines ('%s')", pos, lit)
		}

		ref = append(ref, lit)
	}

	if p.peek() == lexer.TokenIdentifier && p.buf.pos.Line == initialLine {
		_, lit, pos := p.scan()
		return nil, fmt.Errorf("error in %s: unexpected identifier token '%s' at end of reference", pos, lit)
	}

	return &ReferenceNode{Path: ref}, nil
}

func (p *Parser) parseValue() (Node, error) {
	switch p.peek() {
	case lexer.TokenString:
		_, lit, _ := p.scan()
		return &StringLiteral{Value: lit}, nil
	case lexer.TokenNumber:
		return p.parseNumber()
	case lexer.TokenTrue:
		p.scan()
		return &BoolLiteral{Value: true}, nil
	case lexer.TokenFalse:
		p.scan()
		return &BoolLiteral{Value: false}, nil
	case lexer.TokenNull:
		p.scan()
		return &NullLiteral{}, nil
	case lexer.TokenLBracket:
		//parse list
		return p.parseList()
	case lexer.TokenIdentifier:
		return p.parseReference()
	default:
		_, lit, _ := p.scan()
		return nil, errors.New("unexpected token in value: " + lit)
	}
}
