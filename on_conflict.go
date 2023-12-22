package pg

type ConflictTarget interface {
	Update(where ...Condition) ConflictAction
	DoNothing() ConflictAction
}

type ConflictAction interface {
	encodeConflictHandler(b ByteStringWriter, columns []string, args *[]any) error
}

func OnConflict(conflictingColumns ...any) ConflictTarget {
	return onConflict{
		conflictingColumns: Columns(conflictingColumns...),
	}
}

type onConflict struct {
	conflictingColumns Columnar
	targetCondition    Condition
	skip               bool
}

func (c onConflict) Update(where ...Condition) ConflictAction {
	if where != nil {
		c.targetCondition = And(where)
	}

	return c
}

func (c onConflict) DoNothing() ConflictAction {
	c.skip = true
	return c
}

func (c onConflict) encodeConflictHandler(b ByteStringWriter, columns []string, args *[]any) (err error) {
	if c.skip || len(columns) == 0 {
		b.WriteByte('\n')
		b.WriteString("ON CONFLICT DO NOTHING")
		return
	}

	b.WriteByte('\n')
	b.WriteString("ON CONFLICT (")
	c.conflictingColumns.encodeColumnIdentifier(b)
	b.WriteString(") DO UPDATE SET ")

	first := true

	for _, column := range columns {
		if first {
			first = false
		} else {
			b.WriteString(", ")
		}

		writeIdentifier(b, column)
		b.WriteString(" = EXCLUDED.")
		writeIdentifier(b, column)
	}

	if c.targetCondition != nil {
		b.WriteString(" WHERE ")
		c.targetCondition.encodeCondition(b, args)
	}

	return
}
