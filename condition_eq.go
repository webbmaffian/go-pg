package pg

import "strings"

func Eq(column any, value any) Condition {
	c := eq{
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

type eq struct {
	column Columnar
	value  any
}

func (c eq) encodeCondition(b *strings.Builder, args *[]any) {
	c.column.encodeColumnIdentifier(b)

	if c.value == nil {
		b.WriteString(" IS NULL")
	} else {
		b.WriteString(" = ")
		writeParam(b, args, c.value)
	}
}
