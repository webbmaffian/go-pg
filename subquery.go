package pg

import "strings"

var _ Queryable = SubquerySource{}

func Subquery(alias string, query SelectQuery) SubquerySource {
	return SubquerySource{alias, query}
}

type SubquerySource struct {
	alias string
	query SelectQuery
}

func (t SubquerySource) IsZero() bool {
	return t.query.IsZero()
}

func (t SubquerySource) Column(path ...string) AliasedColumnar {
	return subqueryColumn{
		path:  path,
		table: t.alias,
	}
}

func (t SubquerySource) encodeQuery(b *strings.Builder, args *[]any) {
	b.WriteByte('(')
	t.query.encodeQuery(b, args)
	b.WriteByte(')')
	b.WriteString(" AS ")
	writeIdentifier(b, t.alias)
}

type subqueryColumn struct {
	path  []string
	alias string
	table string
}

func (c subqueryColumn) IsZero() bool {
	return c.path == nil
}

func (c subqueryColumn) Alias(alias string) AliasedColumnar {
	c.alias = alias
	return c
}

func (c subqueryColumn) has(column string) bool {
	return c.alias == column
}

func (c subqueryColumn) encodeColumn(b *strings.Builder) {
	c.encode(b)

	if c.alias != "" {
		b.WriteString(" AS ")
		writeIdentifier(b, c.alias)
	}
}

func (c subqueryColumn) encodeColumnIdentifier(b *strings.Builder) {
	if c.alias != "" {
		writeIdentifier(b, c.alias)
	} else {
		c.encode(b)
	}
}

func (c subqueryColumn) encode(b *strings.Builder) {
	l := len(c.path)

	if l > 1 {
		b.WriteByte('(')

		writeIdentifier(b, c.path[0])
		b.WriteString(").")
		c.path = c.path[1:]
	} else if c.table != "" {
		writeIdentifier(b, c.table)
		b.WriteByte('.')
	}

	writeIdentifier(b, c.path...)
}
