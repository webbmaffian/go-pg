package pg

import "strings"

type Columns []Columnar

func (c Columns) encodeColumn(b *strings.Builder) {
	if c == nil {
		return
	}

	c[0].encodeColumn(b)

	for _, v := range c[1:] {
		b.WriteString(", ")
		v.encodeColumn(b)
	}
}

func (c Columns) has(col string) bool {
	for _, column := range c {
		if column.has(col) {
			return true
		}
	}

	return false
}
