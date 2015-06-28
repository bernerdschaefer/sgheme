package main

import (
	"fmt"
	"io"
	"strconv"
	"text/scanner"
	"unicode"
)

var EOF = scanner.EOF

var currentScanner *scmScanner

type scmScanner struct {
	scanner.Scanner
}

func newScanner(r io.Reader) *scmScanner {
	var s scanner.Scanner

	s.Init(r)
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanFloats

	s.IsIdentRune = func(ch rune, i int) bool {
		if i == 0 && unicode.IsDigit(ch) {
			return false
		}

		return ch != '\'' &&
			ch != '"' &&
			ch != '(' &&
			ch != ')' &&
			ch != ';' &&
			!unicode.IsSpace(ch) &&
			unicode.IsPrint(ch)
	}

	return &scmScanner{s}
}

func (s *scmScanner) scanExpression() (x object) {
	var tok = s.Scan()

	for tok != scanner.EOF {
		switch tok {

		case scanner.Int:
			i, _ := strconv.ParseInt(s.TokenText(), 10, 64)
			return scmNumber(i)

		case scanner.Float:
			f, _ := strconv.ParseFloat(s.TokenText(), 64)
			return scmNumber(f)

		case scanner.String:
			str := s.TokenText()
			return scmString(str[1 : len(str)-1])

		case ';':
			s.skipComment()
			tok = s.Scan()

		case '\'':
			return &cell{
				scmSymbol("quote"),
				&cell{s.scanExpression(), NIL},
			}

		case '(':
			return s.scanList()

		case ')':
			return NIL

		case scanner.Ident:
			return scmSymbol(s.TokenText())

		default:
			panic(fmt.Sprintf("unknown token: %s", scanner.TokenString(tok)))

		}
	}

	return EOF
}

func (s *scmScanner) scanList() object {
	var (
		current = &cell{cdr: NIL}
		head    = current
	)

	for e := s.scanExpression(); e != NIL; e = s.scanExpression() {
		new := &cell{car: e, cdr: NIL}
		current.cdr = new
		current = new
	}

	return head.cdr
}

func (s *scmScanner) skipComment() {
	for tok := s.Next(); tok != '\n' && tok != scanner.EOF; tok = s.Next() {
	}
}
