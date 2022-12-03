package pg

import "strings"

func Gt(column any, value any) Condition {
	c := gt{
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

type gt struct {
	column Columnar
	value  any
}

func (c gt) encodeCondition(b *strings.Builder, args *[]any) {
	c.column.encodeColumnIdentifier(b)
	b.WriteString(" > ")
	writeParam(b, args, c.value)
}
