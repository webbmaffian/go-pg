package pg

import "strings"

func RawColumn(raw string) AliasedColumnar {
	return rawColumn{
		raw: raw,
	}
}

type rawColumn struct {
	raw   string
	alias string
}

func (c rawColumn) encodeColumn(b *strings.Builder) {
	b.WriteString(c.raw)

	if c.alias != "" {
		b.WriteString(" AS ")
		writeIdentifier(b, c.alias)
	}
}

func (c rawColumn) encodeColumnIdentifier(b *strings.Builder) {
	if c.alias != "" {
		writeIdentifier(b, c.alias)
	} else {
		b.WriteString(c.raw)
	}
}

func (c rawColumn) Alias(alias string) AliasedColumnar {
	c.alias = alias
	return c
}

func (c rawColumn) has(col string) bool {
	if c.alias != "" {
		return c.alias == col
	}

	return c.raw == col
}
