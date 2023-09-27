package pg

func In(column any, value any) Condition {
	c := in{
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

type in struct {
	column Columnar
	value  any
}

func (c in) IsZero() bool {
	return c.column.IsZero()
}

func (c in) encodeCondition(b ByteStringWriter, args *[]any) {
	c.column.encodeColumnIdentifier(b)

	switch v := c.value.(type) {
	case SelectQuery:
		b.WriteString(" IN (")
		v.encodeQuery(b, args)
		b.WriteByte(')')
	case SubquerySource:
		b.WriteString(" IN (")
		v.query.encodeQuery(b, args)
		b.WriteByte(')')
	default:
		b.WriteString(" = ANY (")
		writeParam(b, args, c.value)
		b.WriteByte(')')
	}

	return
}
