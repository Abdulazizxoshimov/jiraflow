package cql

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

// fieldMap maps CQL field names → SQL columns on the pages table.
var fieldMap = map[string]string{
	"space":       "space_id",
	"title":       "title",
	"type":        "page_type",
	"status":      "status",
	"creator":     "author_id",
	"author":      "author_id",
	"created":     "created_at",
	"lastmodified": "updated_at",
	"updated":     "updated_at",
	"text":        "content_text",
}

// ToSqlizer converts a parsed CQL Query to a squirrel Sqlizer.
// currentUserID is substituted for currentUser() in values.
func ToSqlizer(q *Query, currentUserID string) sq.Sqlizer {
	if len(q.Conditions) == 0 {
		return sq.And{}
	}

	type part struct {
		sql   sq.Sqlizer
		logic string
	}

	var parts []part
	for i, cond := range q.Conditions {
		s := conditionToSql(cond, currentUserID)
		logic := "AND"
		if i < len(q.Logic) {
			logic = q.Logic[i]
		}
		parts = append(parts, part{sql: s, logic: logic})
	}

	result := sq.And{}
	orGroup := sq.Or{}
	for i, p := range parts {
		if p.sql == nil {
			continue
		}
		if i+1 < len(parts) && p.logic == "OR" {
			orGroup = append(orGroup, p.sql)
		} else if len(orGroup) > 0 {
			orGroup = append(orGroup, p.sql)
			result = append(result, orGroup)
			orGroup = sq.Or{}
		} else {
			result = append(result, p.sql)
		}
	}
	if len(orGroup) > 0 {
		result = append(result, orGroup)
	}
	return result
}

// OrderByClause returns the ORDER BY SQL string from parsed query.
func OrderByClause(q *Query) string {
	if len(q.OrderBy) == 0 {
		return ""
	}
	parts := make([]string, 0, len(q.OrderBy))
	for _, o := range q.OrderBy {
		col := resolveField(o.Field)
		if col != "" {
			parts = append(parts, col+" "+o.Direction)
		}
	}
	return strings.Join(parts, ", ")
}

func conditionToSql(c Condition, currentUserID string) sq.Sqlizer {
	vals := resolveValues(c.Values, currentUserID)
	col := resolveField(c.Field)

	// Label subquery
	if c.Field == "label" || c.Field == "labels" {
		return labelSubquerySql(c.Operator, vals)
	}

	if col == "" {
		return nil
	}

	switch c.Operator {
	case "=":
		return sq.Eq{col: vals[0]}
	case "!=":
		return sq.NotEq{col: vals[0]}
	case ">":
		return sq.Gt{col: vals[0]}
	case ">=":
		return sq.GtOrEq{col: vals[0]}
	case "<":
		return sq.Lt{col: vals[0]}
	case "<=":
		return sq.LtOrEq{col: vals[0]}
	case "~":
		return sq.ILike{col: "%" + vals[0] + "%"}
	case "IN":
		return sq.Eq{col: toAny(vals)}
	case "NOT IN":
		return sq.NotEq{col: toAny(vals)}
	case "IS":
		if vals[0] == "EMPTY" || vals[0] == "NULL" {
			return sq.Eq{col: nil}
		}
		return sq.NotEq{col: nil}
	case "IS NOT":
		if vals[0] == "EMPTY" || vals[0] == "NULL" {
			return sq.NotEq{col: nil}
		}
		return sq.Eq{col: nil}
	}
	return nil
}

func resolveField(f string) string {
	return fieldMap[strings.ToLower(f)]
}

func resolveValues(vals []string, currentUserID string) []string {
	out := make([]string, len(vals))
	for i, v := range vals {
		low := strings.ToLower(v)
		if low == "currentuser()" || low == "currentuser" {
			out[i] = currentUserID
		} else {
			out[i] = v
		}
	}
	return out
}

func toAny(ss []string) []any {
	out := make([]any, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}

func labelSubquerySql(op string, vals []string) sq.Sqlizer {
	placeholders := make([]string, len(vals))
	args := make([]any, len(vals))
	for i, v := range vals {
		placeholders[i] = "?"
		args[i] = v
	}
	not := ""
	if op == "NOT IN" || op == "!=" {
		not = "NOT "
	}
	return sq.Expr(
		fmt.Sprintf("id %sIN (SELECT page_id FROM page_labels WHERE label IN (%s))",
			not, strings.Join(placeholders, ",")),
		args...,
	)
}
