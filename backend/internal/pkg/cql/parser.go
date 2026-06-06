package cql

import (
	"fmt"
	"strings"
)

// Condition represents a single field comparison.
type Condition struct {
	Field    string
	Operator string // =, !=, >, >=, <, <=, ~, IN, NOT IN, IS, IS NOT
	Values   []string
}

// OrderClause represents an ORDER BY field + direction.
type OrderClause struct {
	Field     string
	Direction string // ASC | DESC
}

// Query is the parsed CQL result.
type Query struct {
	Conditions []Condition
	Logic      []string // AND | OR between consecutive conditions
	OrderBy    []OrderClause
}

type parser struct {
	tokens []token
	pos    int
}

func (p *parser) peek() token {
	if p.pos >= len(p.tokens) {
		return token{kind: tokEOF}
	}
	return p.tokens[p.pos]
}

func (p *parser) consume() token {
	t := p.peek()
	p.pos++
	return t
}

func (p *parser) expect(k tokenKind) (token, error) {
	t := p.consume()
	if t.kind != k {
		return t, fmt.Errorf("cql: expected token %d got %q", k, t.val)
	}
	return t, nil
}

// Parse parses a CQL string and returns a Query.
func Parse(input string) (*Query, error) {
	tokens := tokenize(input)
	p := &parser{tokens: tokens}
	return p.parseQuery()
}

func (p *parser) parseQuery() (*Query, error) {
	q := &Query{}

	for p.peek().kind != tokEOF && p.peek().kind != tokOrderBy {
		cond, err := p.parseCondition()
		if err != nil {
			return nil, err
		}
		q.Conditions = append(q.Conditions, cond)

		next := p.peek().kind
		if next == tokAnd {
			p.consume()
			q.Logic = append(q.Logic, "AND")
		} else if next == tokOr {
			p.consume()
			q.Logic = append(q.Logic, "OR")
		} else {
			break
		}
	}

	if p.peek().kind == tokOrderBy {
		p.consume()
		for p.peek().kind != tokEOF {
			field := p.consume().val
			dir := "ASC"
			if p.peek().kind == tokAsc {
				p.consume()
			} else if p.peek().kind == tokDesc {
				p.consume()
				dir = "DESC"
			}
			q.OrderBy = append(q.OrderBy, OrderClause{Field: field, Direction: dir})
			if p.peek().kind == tokComma {
				p.consume()
			} else {
				break
			}
		}
	}

	return q, nil
}

func (p *parser) parseCondition() (Condition, error) {
	field, err := p.expect(tokIdent)
	if err != nil {
		return Condition{}, fmt.Errorf("cql: expected field name: %w", err)
	}

	op := p.consume()
	var operator string
	switch op.kind {
	case tokEq:
		operator = "="
	case tokNe:
		operator = "!="
	case tokGt:
		operator = ">"
	case tokGe:
		operator = ">="
	case tokLt:
		operator = "<"
	case tokLe:
		operator = "<="
	case tokTilde:
		operator = "~"
	case tokIn:
		operator = "IN"
	case tokNotIn:
		operator = "NOT IN"
	case tokIs:
		operator = "IS"
	case tokIsNot:
		operator = "IS NOT"
	default:
		return Condition{}, fmt.Errorf("cql: unknown operator %q", op.val)
	}

	var values []string
	if operator == "IN" || operator == "NOT IN" {
		if _, err := p.expect(tokLParen); err != nil {
			return Condition{}, err
		}
		for {
			v := p.consume()
			values = append(values, v.val)
			if p.peek().kind == tokComma {
				p.consume()
			} else {
				break
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return Condition{}, err
		}
	} else if operator == "IS" || operator == "IS NOT" {
		v := p.consume()
		values = []string{strings.ToUpper(v.val)}
	} else {
		v := p.consume()
		values = []string{v.val}
	}

	return Condition{
		Field:    strings.ToLower(field.val),
		Operator: operator,
		Values:   values,
	}, nil
}
