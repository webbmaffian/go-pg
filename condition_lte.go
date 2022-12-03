package pg

import "strings"

func Lte(column any, value any) Condition {
	c := lte{
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

type lte struct {
	column AliasedColumnar
	value  any
}

func (c lte) encodeCondition(b *strings.Builder, args *[]any) {
	c.column.encodeColumnIdentifier(b)
	b.WriteString(" <= ")
	writeParam(b, args, c.value)
}
