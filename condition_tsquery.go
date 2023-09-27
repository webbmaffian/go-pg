package pg

import "strings"

var searchConfig string = "simple"
var prefixSearchSanitizer = strings.NewReplacer(
	`:`, `\:`,
	`*`, `\*`,
	`\`, `\\`,
	`&`, `\&`,
	`|`, `\|`,
)

func SetSearchConfig(conf string) {
	searchConfig = conf
}

func PrefixSearch(vectorColumn any, value string) Condition {
	words := strings.Split(value, " ")
	terms := make([]string, 0, len(words))

	for _, word := range words {
		if word == "" {
			continue
		}

		terms = append(terms, prefixSearchSanitizer.Replace(word)+":*")
	}

	return Search(vectorColumn, strings.Join(terms, " & "))
}

func Search(vectorColumn any, value any) Condition {
	c := tsQuery{
		value: value,
	}

	switch v := vectorColumn.(type) {
	case Columnar:
		c.column = v
	case string:
		c.column = Column(v)
	}

	return c
}

type tsQuery struct {
	column Columnar
	value  any
}

func (c tsQuery) IsZero() bool {
	return c.column.IsZero()
}

func (c tsQuery) encodeCondition(b ByteStringWriter, args *[]any) {
	c.column.encodeColumnIdentifier(b)

	b.WriteString(" @@ to_tsquery('")
	b.WriteString(searchConfig)
	b.WriteString("', ")
	writeParam(b, args, c.value)
	b.WriteByte(')')
}
