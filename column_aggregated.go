package pg

import "strings"

type aggregatedColumn struct {
	function     string
	argsCallback func(b *strings.Builder)
	alias        string
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

func Count(distinctColumn ...string) AliasedColumnar {
	return aggregatedColumn{
		function: "COUNT",
		argsCallback: func(b *strings.Builder) {
			if len(distinctColumn) == 0 {
				b.WriteString("*")
			} else {
				b.WriteString("DISTINCT ")
				writeIdentifier(b, distinctColumn...)
			}
		},
	}
}

func Unnest(columnPath ...string) AliasedColumnar {
	return aggregatedColumn{
		function: "UNNEST",
		argsCallback: func(b *strings.Builder) {
			writeIdentifier(b, columnPath...)
		},
	}
}

func DateTrunc(per string, columnPath ...string) AliasedColumnar {
	return aggregatedColumn{
		function: "date_trunc",
		argsCallback: func(b *strings.Builder) {
			b.WriteByte('\'')
			b.WriteString(per)
			b.WriteString("', ")
			writeIdentifier(b, columnPath...)
		},
	}
}

func Has(columnPath ...string) AliasedColumnar {
	return aggregatedColumn{
		argsCallback: func(b *strings.Builder) {
			writeIdentifier(b, columnPath...)
			b.WriteString(" IS NOT NULL")
		},
	}
}

func Min(columnPath ...string) AliasedColumnar {
	return Aggregate("MIN", columnPath...)
}

func Max(columnPath ...string) AliasedColumnar {
	return Aggregate("MAX", columnPath...)
}

func Sum(columnPath ...string) AliasedColumnar {
	return Aggregate("SUM", columnPath...)
}

func ArrayAgg(columnPath ...string) AliasedColumnar {
	return Aggregate("ARRAY_AGG", columnPath...)
}

func Aggregate(aggFunc string, columnPath ...string) AliasedColumnar {
	return aggregatedColumn{
		function: aggFunc,
		argsCallback: func(b *strings.Builder) {
			writeIdentifier(b, columnPath...)
		},
	}
}
