package pg

func Column(path ...string) AliasedColumnar {
	return column{
		path: path,
	}
}

type column struct {
	path  []string
	alias string
	table *TableSource
}

func (c column) IsZero() bool {
	return c.path == nil
}

func (c column) encodeColumn(b ByteStringWriter) {
	c.encode(b)

	if c.alias != "" {
		b.WriteString(" AS ")
		writeIdentifier(b, c.alias)
	}
}

func (c column) encodeColumnIdentifier(b ByteStringWriter) {
	if c.alias != "" {
		writeIdentifier(b, c.alias)
	} else {
		c.encode(b)
	}
}

func (c column) encode(b ByteStringWriter) {
	l := len(c.path)

	if l > 1 {
		b.WriteByte('(')

		if c.table != nil {
			writeIdentifier(b, c.table.identifier...)
			b.WriteByte('.')
		}

		writeIdentifier(b, c.path[0])
		b.WriteString(").")
		c.path = c.path[1:]
	} else if c.table != nil {
		writeIdentifier(b, c.table.identifier...)
		b.WriteByte('.')
	}

	writeIdentifier(b, c.path...)
}

func (c column) has(col string) bool {
	if c.alias != "" {
		return c.alias == col
	}

	return c.path != nil && c.path[0] == col
}

func (c column) Alias(alias string) AliasedColumnar {
	c.alias = alias
	return c
}
