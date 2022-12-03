package pg

import "strings"

type Join interface {
	encodeJoin(b *strings.Builder, args *[]any)
}

type Joins []Join

func (j Joins) encodeJoin(b *strings.Builder, args *[]any) {
	for _, join := range j {
		join.encodeJoin(b, args)
	}
}

func InnerJoin(table Queryable, condition Condition) Join {
	return join{
		joinType:  "INNER JOIN",
		table:     table,
		condition: condition,
	}
}

func OuterJoin(table Queryable, condition Condition) Join {
	return join{
		joinType:  "OUTER JOIN",
		table:     table,
		condition: condition,
	}
}

func LeftJoin(table Queryable, condition Condition) Join {
	return join{
		joinType:  "LEFT JOIN",
		table:     table,
		condition: condition,
	}
}

func RightJoin(table Queryable, condition Condition) Join {
	return join{
		joinType:  "RIGHT JOIN",
		table:     table,
		condition: condition,
	}
}

type join struct {
	joinType  string
	table     Queryable
	condition Condition
}

func (j join) encodeJoin(b *strings.Builder, args *[]any) {
	b.WriteString(j.joinType)
	b.WriteByte(' ')
	j.table.encodeQuery(b, args)
	b.WriteString(" ON ")
	j.condition.encodeCondition(b, args)
	b.WriteByte('\n')
}
