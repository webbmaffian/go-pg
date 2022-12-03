package pg

import "strings"

func Lt(column any, value any) Condition {
	c := lt{
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

type lt struct {
	column Columnar
	value  any
}

func (c lt) encodeCondition(b *strings.Builder, args *[]any) {
	c.column.encodeColumnIdentifier(b)
	b.WriteString(" < ")
	writeParam(b, args, c.value)
}
