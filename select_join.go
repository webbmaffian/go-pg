package pg

import "strings"

type join interface {
	encodeJoin(b *strings.Builder, args *[]any)
}

type Joins []join

func (j Joins) encodeJoin(b *strings.Builder, args *[]any) {
	for _, join := range j {
		join.encodeJoin(b, args)
	}
}

type InnerJoin struct {
	Table     queryable
	Condition Condition
}

func (j InnerJoin) encodeJoin(b *strings.Builder, args *[]any) {
	b.WriteString("INNER JOIN ")
	j.Table.encodeQuery(b, args)
	b.WriteString(" ON ")
	j.Condition.encodeCondition(b, args)
	b.WriteByte('\n')
}

type OuterJoin InnerJoin

func (j OuterJoin) encodeJoin(b *strings.Builder, args *[]any) {
	b.WriteString("OUTER JOIN ")
	j.Table.encodeQuery(b, args)
	b.WriteString(" ON ")
	j.Condition.encodeCondition(b, args)
	b.WriteByte('\n')
}

type LeftJoin InnerJoin

func (j LeftJoin) encodeJoin(b *strings.Builder, args *[]any) {
	b.WriteString("LEFT JOIN ")
	j.Table.encodeQuery(b, args)
	b.WriteString(" ON ")
	j.Condition.encodeCondition(b, args)
	b.WriteByte('\n')
}

type RightJoin InnerJoin

func (j RightJoin) encodeJoin(b *strings.Builder, args *[]any) {
	b.WriteString("RIGHT JOIN ")
	j.Table.encodeQuery(b, args)
	b.WriteString(" ON ")
	j.Condition.encodeCondition(b, args)
	b.WriteByte('\n')
}
