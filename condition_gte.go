package pg

import "strings"

func Gte(column any, value any) Condition {
	c := gte{
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

type gte struct {
	column Columnar
	value  any
}

func (c gte) IsZero() bool {
	return c.column.IsZero()
}

func (c gte) encodeCondition(b *strings.Builder, args *[]any) {
	c.column.encodeColumnIdentifier(b)
	b.WriteString(" >= ")
	writeParam(b, args, c.value)
}
