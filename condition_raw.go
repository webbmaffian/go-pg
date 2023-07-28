package pg

import (
	"bytes"
	"strings"
)

func Raw(str string, params ...any) Condition {
	return raw{
		String: str,
		Params: params,
	}
}

type raw struct {
	String string
	Params []any
}

func (c raw) IsZero() bool {
	return c.String == ""
}

func (c raw) encodeCondition(b *strings.Builder, args *[]any) {
	if len(c.Params) == 0 {
		b.WriteString(c.String)
		return
	}

	var prev int
	str := []byte(c.String)

	for _, param := range c.Params {
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
