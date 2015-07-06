package main

import (
	"io"
	"strconv"
	"text/scanner"
	"unicode"
)

const EOF = scanner.EOF

var currentScanner *scmScanner

type scmScanner struct {
	scanner.Scanner

	expectedCloses int
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

	return &scmScanner{Scanner: s}
}

func (s *scmScanner) read() object {
	var tok = s.Scan()

	for tok != EOF {
		switch tok {

		case scanner.Ident:
			return scmSymbol(s.TokenText())

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
			return s.readQuoted()

		case '(':
			return s.readList()

		case ')':
			if s.expectedCloses > 0 {
				return ')'
			}

			fallthrough

		default:
			return raiseError(
				"Syntax error, invalid token",
				scanner.TokenString(tok),
			)

		}
	}

	return EOF
}

func (s *scmScanner) readList() object {
	var (
		current = &cell{cdr: NIL}
		head    = current
	)

	for {
		s.expectedCloses++
		e := s.read()
		s.expectedCloses--

		switch e {
		case EOF:
			return EOF
		case ')':
			return head.cdr
		}

		new := &cell{car: e, cdr: NIL}
		current.cdr = new
		current = new
	}
}

func (s *scmScanner) readQuoted() object {
	e := s.read()
	if e == EOF {
		return EOF
	}

	return &cell{
		scmSymbol("quote"),
		&cell{e, NIL},
	}
}

func (s *scmScanner) skipComment() {
	for tok := s.Next(); tok != '\n' && tok != EOF; tok = s.Next() {
	}
}
