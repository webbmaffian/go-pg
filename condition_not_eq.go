package pg

func NotEq(column any, value any) Condition {
	c := notEq{
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

type notEq struct {
	column Columnar
	value  any
}

func (c notEq) IsZero() bool {
	return c.column.IsZero()
}

func (c notEq) encodeCondition(b ByteStringWriter, args *[]any) {
	c.column.encodeColumnIdentifier(b)

	if c.value == nil {
		b.WriteString(" IS NOT NULL")
	} else {
		b.WriteString(" != ")
		writeParam(b, args, c.value)
	}
}
