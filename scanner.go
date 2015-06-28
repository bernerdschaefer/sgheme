package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"text/scanner"
	"unicode"
)

var currentScanner *scanner.Scanner

type scanner scanner.Scanner

func newScanner(r io.Reader) *scanner.Scanner {
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

	return &s
}

func scanExpression() (x object) {
	var tok = currentScanner.Scan()

	for tok != scanner.EOF {
		switch tok {

		case scanner.Int:
			i, _ := strconv.ParseInt(currentScanner.TokenText(), 10, 64)
			return scmNumber(i)

		case scanner.Float:
			f, _ := strconv.ParseFloat(currentScanner.TokenText(), 64)
			return scmNumber(f)

		case scanner.String:
			str := currentScanner.TokenText()
			return scmString(str[1 : len(str)-1])

		case ';':
			skipComment()
			tok = currentScanner.Scan()

		case '\'':
			return &cell{
				scmSymbol("quote"),
				&cell{scanExpression(), NIL},
			}

		case '(':
			return scanList()

		case ')':
			return NIL

		case scanner.Ident:
			return scmSymbol(currentScanner.TokenText())

		default:
			panic(fmt.Sprintf("unknown token: %s", scanner.TokenString(tok)))

		}
	}

	os.Exit(0)
	panic("unreachable")
}

func scanList() object {
	var (
		current = &cell{cdr: NIL}
		head    = current
	)

	for e := scanExpression(); e != NIL; e = scanExpression() {
		new := &cell{car: e, cdr: NIL}
		current.cdr = new
		current = new
	}

	return head.cdr
}

func skipComment() {
	for tok := currentScanner.Next(); tok != '\n' && tok != scanner.EOF; tok = currentScanner.Next() {
	}
}
