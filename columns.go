package pg

func Columns(cols ...any) MultiColumnar {
	realCols := make(columns, len(cols))

	for i := range cols {
		switch c := cols[i].(type) {

		case Columnar:
			realCols[i] = c

		case string:
			realCols[i] = Column(c)
		}

	}

	return realCols
}

func AllocateColumns(capacity int) MultiColumnar {
	return make(columns, 0, capacity)
}

type columns []Columnar

func (c columns) IsZero() bool {
	for i := range c {
		if !c[i].IsZero() {
			return false
		}
	}

	return true
}

func (c columns) Append(cols ...Columnar) MultiColumnar {
	return append(c, cols...)
}

func (c columns) encodeColumn(b ByteStringWriter) {
	if c == nil {
		return
	}

	c[0].encodeColumn(b)

	for _, v := range c[1:] {
		b.WriteString(", ")
		v.encodeColumn(b)
	}
}

func (c columns) encodeColumnIdentifier(b ByteStringWriter) {
	if c == nil {
		return
	}

	c[0].encodeColumnIdentifier(b)

	for _, v := range c[1:] {
		b.WriteString(", ")
		v.encodeColumnIdentifier(b)
	}
}

func (c columns) has(col string) bool {
	for _, column := range c {
		if column.has(col) {
			return true
		}
	}

	return false
}
