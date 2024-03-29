package pg

func NotIn(column any, value any) Condition {
	c := notIn{
		value: value,
	}

	switch v := column.(type) {
	case Columnar:
		c.column = v
	case string:
		c.column = Column(v)
	}

	return c
}

type notIn struct {
	column Columnar
	value  any
}

func (c notIn) IsZero() bool {
	return c.column.IsZero()
}

func (c notIn) encodeCondition(b ByteStringWriter, args *[]any) {
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
