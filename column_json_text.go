package pg

func JsonColumnText(path ...string) AliasedColumnar {
	return jsonColumnText{
		path: path,
	}
}

type jsonColumnText struct {
	path  []string
	alias string
	table *TableSource
}

func (c jsonColumnText) IsZero() bool {
	return c.path == nil
}

func (c jsonColumnText) encodeColumn(b ByteStringWriter) {
	c.encode(b)

	if c.alias != "" {
		b.WriteString(" AS ")
		writeIdentifier(b, c.alias)
	}
}

func (c jsonColumnText) encodeColumnIdentifier(b ByteStringWriter) {
	if c.alias != "" {
		writeIdentifier(b, c.alias)
	} else {
		c.encode(b)
	}
}

func (c jsonColumnText) encode(b ByteStringWriter) {
	if c.path == nil {
		return
	}

	if c.table != nil {
		writeIdentifier(b, c.table.identifier...)
		b.WriteByte('.')
	}

	writeIdentifier(b, c.path[0])

	for _, col := range c.path[1:] {
		b.WriteString("->>")
		writeString(b, col)
	}
}

func (c jsonColumnText) has(col string) bool {
	if c.alias != "" {
		return c.alias == col
	}

	return c.path != nil && c.path[0] == col
}

func (c jsonColumnText) Alias(alias string) AliasedColumnar {
	c.alias = alias
	return c
}
