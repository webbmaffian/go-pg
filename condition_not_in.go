package pg

import "strings"

func NotIn(column any, value any) Condition {
	c := notIn{
		value: value,
	}

	switch v := column.(type) {
	case AliasedColumnar:
		c.column = v
	case string:
		c.column = Column(v)
	}

	return c
}

type notIn struct {
	column AliasedColumnar
	value  any
}

func (c notIn) encodeCondition(b *strings.Builder, args *[]any) {
	c.column.encodeColumnIdentifier(b)

	switch v := c.value.(type) {
	case SelectQuery:
		b.WriteString(" NOT IN (")
		v.encodeQuery(b, args)
		b.WriteByte(')')
	case SubquerySource:
		b.WriteString(" NOT IN (")
		v.query.encodeQuery(b, args)
		b.WriteByte(')')
	default:
		b.WriteString(" != ANY (")
		writeParam(b, args, c.value)
		b.WriteByte(')')
	}

	return
}
