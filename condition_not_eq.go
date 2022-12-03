package pg

import "strings"

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

func (c notEq) encodeCondition(b *strings.Builder, args *[]any) {
	c.column.encodeColumnIdentifier(b)

	if c.value == nil {
		b.WriteString(" IS NOT NULL")
	} else {
		b.WriteString(" != ")
		writeParam(b, args, c.value)
	}
}
