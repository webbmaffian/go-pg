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

func (t SubquerySource) encodeQuery(b *strings.Builder, args *[]any) {
	b.WriteByte('(')
	t.query.encodeQuery(b, args)
	b.WriteByte(')')
	b.WriteString(" AS ")
	writeIdentifier(b, t.alias)
}
