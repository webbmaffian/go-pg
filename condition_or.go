package pg

import "strings"

var _ Condition = Or{}

type Or []Condition

func (c Or) IsZero() bool {
	return c == nil
}

func (c Or) encodeCondition(b *strings.Builder, args *[]any) {
	if c == nil {
		return
	}

	b.WriteByte('(')
	c[0].encodeCondition(b, args)

	for _, cond := range c[1:] {
		b.WriteString(" OR ")
		cond.encodeCondition(b, args)
	}

	b.WriteByte(')')
}
