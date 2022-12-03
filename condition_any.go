package pg

import "strings"

func Any(column any, value any) Condition {
	c := some{
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

func Some(column any, value any) Condition {
	return Any(column, value)
}

type some struct {
	column AliasedColumnar
	value  any
}

func (c some) encodeCondition(b *strings.Builder, args *[]any) {
	writeParam(b, args, c.value)
	b.WriteString(" = ANY (")
	c.column.encodeColumnIdentifier(b)
	b.WriteByte(')')

	return
}
