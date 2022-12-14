package pg

import "strings"

type aggregatedColumn struct {
	function     string
	argsCallback func(b *strings.Builder)
	alias        string
}

func (c aggregatedColumn) IsZero() bool {
	return c.function == "" && c.argsCallback == nil && c.alias == ""
}

func (c aggregatedColumn) encodeColumn(b *strings.Builder) {
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

func (c aggregatedColumn) encodeColumnIdentifier(b *strings.Builder) {
	writeIdentifier(b, c.alias)
}

func (c aggregatedColumn) has(column string) bool {
	return c.alias == column
}

func Count(distinctColumn ...any) AliasedColumnar {
	return aggregatedColumn{
		function: "COUNT",
		argsCallback: func(b *strings.Builder) {
			if len(distinctColumn) == 0 {
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
		argsCallback: func(b *strings.Builder) {
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
		argsCallback: func(b *strings.Builder) {
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

func ArrayAgg(column string) AliasedColumnar {
	return Aggregate("ARRAY_AGG", column)
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
