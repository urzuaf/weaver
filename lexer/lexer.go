package lexer

import (
	"fmt"
	"io"
	"strings"
	"text/scanner"
)

type Token int

const (
	TokenError Token = iota
	TokenEOF
	TokenIdentifier
	TokenString
	TokenNumber
	TokenTrue
	TokenFalse
	TokenNull
	TokenLBrace   // {
	TokenRBrace   // }
	TokenEqual    // =
	TokenColon    // :
	TokenLBracket // [
	TokenRBracket // ]
	TokenComma    // ,
	TokenDot      // .

)

var charTokens = map[rune]Token{
	'{': TokenLBrace,
	'}': TokenRBrace,
	'=': TokenEqual,
	':': TokenColon,
	'[': TokenLBracket,
	']': TokenRBracket,
	',': TokenComma,
	'.': TokenDot,
}

type Lexer struct {
	s       scanner.Scanner
	verbose bool
}

func NewLexer(src io.Reader, verbose bool) *Lexer {
	var s scanner.Scanner
	s.Init(src)
	return &Lexer{s: s}
}

func (l *Lexer) Next() (tok Token, lit string, pos scanner.Position) {
	for {
		stok := l.s.Scan()
		pos = l.s.Position
		lit = l.s.TokenText()

		if l.verbose {
			fmt.Println("Lexer token:", stok, "lit:", lit, "pos:", pos)
		}
		if stok == '#' {
			for l.s.Peek() != '\n' && l.s.Peek() != scanner.EOF {
				l.s.Next()
			}
			continue
		}

		switch stok {
		case scanner.EOF:
			return TokenEOF, "", pos
		case scanner.Ident:
			switch strings.ToLower(lit) {
			case "true":
				return TokenTrue, lit, pos
			case "false":
				return TokenFalse, lit, pos
			case "null":
				return TokenNull, lit, pos
			}
			return TokenIdentifier, lit, pos
		case scanner.String:
			return TokenString, lit[1 : len(lit)-1], pos // remove quotes
		case scanner.Int, scanner.Float:
			return TokenNumber, lit, pos
		default:
			if token, ok := charTokens[stok]; ok {
				return token, lit, pos
			}
			return TokenError, lit, pos
		}
	}
}

func (t *Token) ToString() string {
	return fmt.Sprint(*t)
}
