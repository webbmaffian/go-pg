package pg

import "strconv"

type aggregatedColumn struct {
	function     string
	argsCallback func(b ByteStringWriter)
	alias        string
}

func (c aggregatedColumn) IsZero() bool {
	return c.function == "" && c.argsCallback == nil && c.alias == ""
}

func (c aggregatedColumn) encodeColumn(b ByteStringWriter) {
	if c.function != "" {
		b.WriteString(c.function)
	}

	if c.argsCallback != nil {
		b.WriteByte('(')
		c.argsCallback(b)
		b.WriteByte(')')
	}

	if c.alias != "" {
		b.WriteString(" AS ")
		writeIdentifier(b, c.alias)
	}
}

func (c aggregatedColumn) Alias(alias string) AliasedColumnar {
	c.alias = alias
	return c
}

func (c aggregatedColumn) encodeColumnIdentifier(b ByteStringWriter) {
	writeIdentifier(b, c.alias)
}

func (c aggregatedColumn) has(column string) bool {
	return c.alias == column
}

func Count(distinctColumn ...any) AliasedColumnar {
	return aggregatedColumn{
		function: "COUNT",
		argsCallback: func(b ByteStringWriter) {
			if len(distinctColumn) == 0 || distinctColumn[0] == nil {
				b.WriteString("*")
			} else {
				var col Columnar

				switch v := distinctColumn[0].(type) {
				case Columnar:
					col = v
				case string:
					col = Column(v)
				}
				b.WriteString("DISTINCT ")
				col.encodeColumnIdentifier(b)
			}
		},
	}
}

func Unnest(column any) AliasedColumnar {
	return Aggregate("UNNEST", column)
}

func DateTrunc(per string, column any) AliasedColumnar {
	var col Columnar

	switch v := column.(type) {
	case Columnar:
		col = v
	case string:
		col = Column(v)
	}

	return aggregatedColumn{
		function: "date_trunc",
		argsCallback: func(b ByteStringWriter) {
			b.WriteByte('\'')
			b.WriteString(per)
			b.WriteString("', ")
			col.encodeColumnIdentifier(b)
		},
	}
}

func Has(column any) AliasedColumnar {
	var col Columnar

	switch v := column.(type) {
	case Columnar:
		col = v
	case string:
		col = Column(v)
	}

	return aggregatedColumn{
		argsCallback: func(b ByteStringWriter) {
			col.encodeColumnIdentifier(b)
			b.WriteString(" IS NOT NULL")
		},
	}
}

func Min(column any) AliasedColumnar {
	return Aggregate("MIN", column)
}

func Max(column any) AliasedColumnar {
	return Aggregate("MAX", column)
}

func Sum(column any) AliasedColumnar {
	return Aggregate("SUM", column)
}

func ArrayAgg(column any, orderBy ...OrderByColumnar) AliasedColumnar {
	var col Columnar

	switch v := column.(type) {
	case Columnar:
		col = v
	case string:
		col = Column(v)
	}

	return aggregatedColumn{
		function: "ARRAY_AGG",
		argsCallback: func(b ByteStringWriter) {
			col.encodeColumnIdentifier(b)

			if orderBy != nil && orderBy[0] != nil {
				b.WriteString(" ORDER BY ")
				orderBy[0].encodeOrderBy(b)
			}
		},
	}
}

func Aggregate(aggFunc string, column any) AliasedColumnar {
	var col Columnar

	switch v := column.(type) {
	case Columnar:
		col = v
	case string:
		col = Column(v)
	}

	return aggregatedColumn{
		function:     aggFunc,
		argsCallback: col.encodeColumnIdentifier,
	}
}

func Coalesce(columns ...any) AliasedColumnar {
	return aggregatedColumn{
		function: "COALESCE",
		argsCallback: func(b ByteStringWriter) {
			for i := range columns {
				var col Columnar

				switch v := columns[i].(type) {
				case Columnar:
					col = v
				case string:
					col = Column(v)
				case int:
					col = Raw(strconv.Itoa(v))
				}

				if i != 0 {
					b.WriteString(", ")
				}

				col.encodeColumn(b)
			}
		},
	}
}
