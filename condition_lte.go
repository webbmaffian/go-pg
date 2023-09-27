package pg

func Lte(column any, value any) Condition {
	c := lte{
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

type lte struct {
	column Columnar
	value  any
}

func (c lte) IsZero() bool {
	return c.column.IsZero()
}

func (c lte) encodeCondition(b ByteStringWriter, args *[]any) {
	c.column.encodeColumnIdentifier(b)
	b.WriteString(" <= ")
	writeParam(b, args, c.value)
}
