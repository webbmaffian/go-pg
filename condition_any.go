package pg

func Any(column any, value any) Condition {
	c := some{
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

func Some(column any, value any) Condition {
	return Any(column, value)
}

type some struct {
	column Columnar
	value  any
}

func (c some) IsZero() bool {
	return c.column.IsZero()
}

func (c some) encodeCondition(b ByteStringWriter, args *[]any) {
	writeParam(b, args, c.value)
	b.WriteString(" = ANY (")
	c.column.encodeColumnIdentifier(b)
	b.WriteByte(')')
}
