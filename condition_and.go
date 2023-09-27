package pg

var _ Condition = And{}

type And []Condition

func (c And) IsZero() bool {
	return c == nil
}

func (c And) encodeCondition(b ByteStringWriter, args *[]any) {
	if c == nil {
		return
	}

	b.WriteByte('(')
	c[0].encodeCondition(b, args)

	for _, cond := range c[1:] {
		b.WriteString(" AND ")
		cond.encodeCondition(b, args)
	}

	b.WriteByte(')')
}
