package pg

import "strings"

func Columns(cols ...any) Columnar {
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

type columns []Columnar

func (c columns) encodeColumn(b *strings.Builder) {
	if c == nil {
		return
	}

	c[0].encodeColumn(b)

	for _, v := range c[1:] {
		b.WriteString(", ")
		v.encodeColumn(b)
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
