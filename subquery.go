package pg

import "strings"

func Subquery(alias string, query SelectQuery) SubquerySource {
	return SubquerySource{alias, query}
}

type SubquerySource struct {
	alias string
	query SelectQuery
}

func (t SubquerySource) encodeQuery(b *strings.Builder, args *[]any) {
	b.WriteByte('(')
	t.query.encodeQuery(b, args)
	b.WriteByte(')')
	b.WriteString(" AS ")
	writeIdentifier(b, t.alias)
}
