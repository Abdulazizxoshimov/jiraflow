package cql

import (
	"strings"
	"unicode"
)

type tokenKind int

const (
	tokEOF tokenKind = iota
	tokIdent
	tokString
	tokEq
	tokNe
	tokGt
	tokGe
	tokLt
	tokLe
	tokTilde
	tokLParen
	tokRParen
	tokComma
	tokAnd
	tokOr
	tokNot
	tokIn
	tokNotIn
	tokIs
	tokIsNot
	tokOrderBy
	tokAsc
	tokDesc
	tokEmpty
	tokNull
)

type token struct {
	kind tokenKind
	val  string
}

type lexer struct {
	input []rune
	pos   int
}

func newLexer(s string) *lexer { return &lexer{input: []rune(s)} }

func (l *lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *lexer) next() rune {
	r := l.peek()
	l.pos++
	return r
}

func (l *lexer) skipWS() {
	for l.pos < len(l.input) && unicode.IsSpace(l.input[l.pos]) {
		l.pos++
	}
}

func (l *lexer) readString(quote rune) string {
	var sb strings.Builder
	for l.pos < len(l.input) {
		r := l.next()
		if r == quote {
			break
		}
		if r == '\\' && l.pos < len(l.input) {
			sb.WriteRune(l.next())
			continue
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

func (l *lexer) readIdent() string {
	var sb strings.Builder
	for l.pos < len(l.input) {
		r := l.peek()
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' || r == '-' || r == '~' {
			sb.WriteRune(r)
			l.pos++
		} else {
			break
		}
	}
	return sb.String()
}

var keywords = map[string]tokenKind{
	"and":     tokAnd,
	"or":      tokOr,
	"not":     tokNot,
	"in":      tokIn,
	"is":      tokIs,
	"asc":     tokAsc,
	"desc":    tokDesc,
	"empty":   tokEmpty,
	"null":    tokNull,
	"orderby": tokOrderBy,
}

func (l *lexer) Scan() token {
	l.skipWS()
	if l.pos >= len(l.input) {
		return token{kind: tokEOF}
	}

	r := l.peek()

	if r == '"' || r == '\'' {
		l.pos++
		return token{kind: tokString, val: l.readString(r)}
	}

	if l.pos+1 < len(l.input) {
		two := string(l.input[l.pos : l.pos+2])
		switch two {
		case "!=":
			l.pos += 2
			return token{kind: tokNe, val: "!="}
		case ">=":
			l.pos += 2
			return token{kind: tokGe, val: ">="}
		case "<=":
			l.pos += 2
			return token{kind: tokLe, val: "<="}
		}
		rest := strings.ToLower(string(l.input[l.pos:]))
		if strings.HasPrefix(rest, "order by") {
			l.pos += 8
			return token{kind: tokOrderBy, val: "ORDER BY"}
		}
		if strings.HasPrefix(rest, "not in") {
			l.pos += 6
			return token{kind: tokNotIn, val: "NOT IN"}
		}
		if strings.HasPrefix(rest, "is not") {
			l.pos += 6
			return token{kind: tokIsNot, val: "IS NOT"}
		}
	}

	l.pos++
	switch r {
	case '=':
		return token{kind: tokEq, val: "="}
	case '>':
		return token{kind: tokGt, val: ">"}
	case '<':
		return token{kind: tokLt, val: "<"}
	case '~':
		return token{kind: tokTilde, val: "~"}
	case '(':
		return token{kind: tokLParen, val: "("}
	case ')':
		return token{kind: tokRParen, val: ")"}
	case ',':
		return token{kind: tokComma, val: ","}
	}

	l.pos--
	if unicode.IsLetter(r) || r == '_' {
		ident := l.readIdent()
		lower := strings.ToLower(ident)
		if kind, ok := keywords[lower]; ok {
			return token{kind: kind, val: ident}
		}
		return token{kind: tokIdent, val: ident}
	}

	if unicode.IsDigit(r) || r == '-' {
		return token{kind: tokString, val: l.readIdent()}
	}

	l.pos++
	return token{kind: tokIdent, val: string(r)}
}

func tokenize(input string) []token {
	l := newLexer(input)
	var tokens []token
	for {
		t := l.Scan()
		tokens = append(tokens, t)
		if t.kind == tokEOF {
			break
		}
	}
	return tokens
}
