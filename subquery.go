package pg

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

func (t SubquerySource) Alias() string {
	return t.alias
}

func (t SubquerySource) Query() Queryable {
	return &t.query
}

func (t SubquerySource) encodeQuery(b ByteStringWriter, args *[]any) {
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

func (c subqueryColumn) encodeColumn(b ByteStringWriter) {
	c.encode(b)

	if c.alias != "" {
		b.WriteString(" AS ")
		writeIdentifier(b, c.alias)
	}
}

func (c subqueryColumn) encodeColumnIdentifier(b ByteStringWriter) {
	if c.alias != "" {
		writeIdentifier(b, c.alias)
	} else {
		c.encode(b)
	}
}

func (c subqueryColumn) encode(b ByteStringWriter) {
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
