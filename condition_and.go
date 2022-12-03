package pg

import "strings"

type And []Condition

func (c And) encodeCondition(b *strings.Builder, args *[]any) {
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
