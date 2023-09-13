package pg

import (
	"bytes"
	"strings"
)

// Any raw part of the query. Params will only be used for conditions.
func Raw(str string, params ...any) RawData {
	return raw{
		raw:    str,
		params: params,
	}
}

// Deprecated: Use Raw() instead
func RawColumn(str string) RawData {
	return raw{
		raw: str,
	}
}

type raw struct {
	params []any
	raw    string
	alias  string
}

func (c raw) IsZero() bool {
	return c.raw == ""
}

func (c raw) encodeCondition(b *strings.Builder, args *[]any) {
	if len(c.params) == 0 {
		b.WriteString(c.raw)
		return
	}

	var prev int
	str := []byte(c.raw)

	for _, param := range c.params {
		cur := bytes.IndexByte(str[prev:], '?')

		if cur == -1 {
			break
		}

		b.Write(str[prev:cur])
		writeParam(b, args, param)

		prev = cur + 1
	}

	b.Write(str[prev:])
}

func (c raw) encodeColumn(b *strings.Builder) {
	b.WriteString(c.raw)

	if c.alias != "" {
		b.WriteString(" AS ")
		writeIdentifier(b, c.alias)
	}
}

func (c raw) encodeColumnIdentifier(b *strings.Builder) {
	if c.alias != "" {
		writeIdentifier(b, c.alias)
	} else {
		b.WriteString(c.raw)
	}
}

func (c raw) Alias(alias string) RawData {
	c.alias = alias
	return c
}

func (c raw) has(col string) bool {
	if c.alias != "" {
		return c.alias == col
	}

	return c.raw == col
}

func (c raw) encodeQuery(b *strings.Builder, args *[]any) {
	b.WriteString(c.raw)

	if c.alias != "" {
		b.WriteString(" AS ")
		b.WriteString(c.alias)
	}
}
