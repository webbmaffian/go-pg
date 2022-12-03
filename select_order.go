package pg

import "strings"

type orderBy interface {
	encodeOrderBy(b *strings.Builder)
}

type OrderBy []orderBy

func (o OrderBy) encodeOrderBy(b *strings.Builder) {
	if len(o) == 0 {
		return
	}

	o[0].encodeOrderBy(b)

	for _, v := range o[1:] {
		b.WriteString(", ")
		v.encodeOrderBy(b)
	}
}

type Asc []AliasedColumnar

func (o Asc) encodeOrderBy(b *strings.Builder) {
	if len(o) == 0 {
		return
	}

	o[0].encodeColumnIdentifier(b)
	b.WriteString(" ASC")

	for _, v := range o[1:] {
		b.WriteString(", ")
		v.encodeColumnIdentifier(b)
		b.WriteString(" ASC")
	}
}

type Desc []AliasedColumnar

func (o Desc) encodeOrderBy(b *strings.Builder) {
	if len(o) == 0 {
		return
	}

	o[0].encodeColumnIdentifier(b)
	b.WriteString(" DESC")

	for _, v := range o[1:] {
		b.WriteString(", ")
		v.encodeColumnIdentifier(b)
		b.WriteString(" DESC")
	}
}
