package pg

import (
	"strings"
)

type OnConflictUpdate []string

func (conflictingColumns OnConflictUpdate) run(b *strings.Builder, columns []string) (err error) {
	if len(columns) == 0 {
		return
	}

	b.WriteByte('\n')
	b.WriteString("ON CONFLICT (")
	writeIdentifier(b, conflictingColumns[0])

	for _, column := range conflictingColumns[1:] {
		b.WriteString(", ")
		writeIdentifier(b, column)
	}

	b.WriteString(") DO UPDATE SET ")

	first := true

	for _, column := range columns {
		if first {
			first = false
		} else {
			b.WriteString(", ")
		}

		writeIdentifier(b, column)
		b.WriteString(" = excluded.")
		writeIdentifier(b, column)
	}

	return
}
