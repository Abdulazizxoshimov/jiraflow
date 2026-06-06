package jql

import (
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
)

// fieldMap maps JQL field names → SQL column names.
var fieldMap = map[string]string{
	"project":      "project_id",
	"type":         "type",
	"issuetype":    "type",
	"status":       "status_id",
	"priority":     "priority",
	"assignee":     "assignee_id",
	"reporter":     "reporter_id",
	"sprint":       "sprint_id",
	"parent":       "parent_id",
	"resolution":   "resolution",
	"created":      "created_at",
	"updated":      "updated_at",
	"duedate":      "due_date",
	"due":          "due_date",
	"summary":      "title",
	"text":         "title",
	"description":  "description",
}

// ToSqlizer converts a parsed Query to a squirrel Sqlizer (AND of all conditions).
// currentUserID is substituted for currentUser() function in values.
func ToSqlizer(q *Query, currentUserID string) sq.Sqlizer {
	if len(q.Conditions) == 0 {
		return sq.And{}
	}

	type part struct {
		sql  sq.Sqlizer
		logic string // AND | OR (what follows this part)
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

	// Build AND groups (OR has lower precedence but we keep it simple: left-to-right)
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

	// Special subquery fields
	switch c.Field {
	case "labels", "label":
		return subquerySql("issue_labels", "label_id", c.Operator, vals)
	case "component", "components":
		return subquerySql("issue_components", "component_id", c.Operator, vals)
	case "fixversion", "fixversions", "fix_version":
		return versionSubquerySql("fix", c.Operator, vals)
	case "affectsversion", "affectsversions", "affects_version":
		return versionSubquerySql("affects", c.Operator, vals)
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
		anys := toAny(vals)
		return sq.Eq{col: anys}
	case "NOT IN":
		anys := toAny(vals)
		return sq.NotEq{col: anys}
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
	now := time.Now().UTC()
	out := make([]string, len(vals))
	for i, v := range vals {
		out[i] = resolveValue(v, currentUserID, now)
	}
	return out
}

// resolveValue converts a single JQL value token to its SQL equivalent.
// Handles currentUser(), date functions, and relative date offsets like -1d, +2w.
func resolveValue(v, currentUserID string, now time.Time) string {
	low := strings.ToLower(strings.TrimSpace(v))

	switch low {
	case "currentuser()", "currentuser":
		return currentUserID
	case "startofday()":
		return startOf(now, "day").Format(time.RFC3339)
	case "endofday()":
		return endOf(now, "day").Format(time.RFC3339)
	case "startofweek()":
		return startOf(now, "week").Format(time.RFC3339)
	case "endofweek()":
		return endOf(now, "week").Format(time.RFC3339)
	case "startofmonth()":
		return startOf(now, "month").Format(time.RFC3339)
	case "endofmonth()":
		return endOf(now, "month").Format(time.RFC3339)
	case "startofyear()":
		return startOf(now, "year").Format(time.RFC3339)
	case "endofyear()":
		return endOf(now, "year").Format(time.RFC3339)
	}

	// Relative date offsets: -1d, +2w, -3m, +1y
	if t, ok := parseRelativeDate(low, now); ok {
		return t.Format(time.RFC3339)
	}

	return v
}

func startOf(now time.Time, unit string) time.Time {
	switch unit {
	case "day":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // treat Sunday as 7
		}
		return time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, time.UTC)
	case "month":
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	case "year":
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return now
}

func endOf(now time.Time, unit string) time.Time {
	switch unit {
	case "day":
		return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		return time.Date(now.Year(), now.Month(), now.Day()+(7-weekday), 23, 59, 59, 0, time.UTC)
	case "month":
		first := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		return first.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	case "year":
		return time.Date(now.Year(), 12, 31, 23, 59, 59, 0, time.UTC)
	}
	return now
}

// parseRelativeDate handles -1d, +2w, -3m, +1y offsets.
func parseRelativeDate(s string, now time.Time) (time.Time, bool) {
	if len(s) < 2 {
		return now, false
	}
	sign := 1
	rest := s
	if s[0] == '-' {
		sign = -1
		rest = s[1:]
	} else if s[0] == '+' {
		rest = s[1:]
	} else {
		return now, false
	}
	if len(rest) < 2 {
		return now, false
	}
	unit := rest[len(rest)-1]
	var n int
	if _, err := fmt.Sscanf(rest[:len(rest)-1], "%d", &n); err != nil {
		return now, false
	}
	n *= sign
	switch unit {
	case 'd':
		return now.AddDate(0, 0, n), true
	case 'w':
		return now.AddDate(0, 0, n*7), true
	case 'm':
		return now.AddDate(0, n, 0), true
	case 'y':
		return now.AddDate(n, 0, 0), true
	}
	return now, false
}

func toAny(ss []string) []any {
	out := make([]any, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}

func subquerySql(table, idCol, op string, vals []string) sq.Sqlizer {
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
		fmt.Sprintf("id %sIN (SELECT issue_id FROM %s WHERE %s IN (%s))",
			not, table, idCol, strings.Join(placeholders, ",")),
		args...,
	)
}

func versionSubquerySql(vType, op string, vals []string) sq.Sqlizer {
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
		fmt.Sprintf("id %sIN (SELECT issue_id FROM issue_versions WHERE version_id IN (%s) AND version_type = '%s')",
			not, strings.Join(placeholders, ","), vType),
		args...,
	)
}
