package pg

import (
	"database/sql/driver"
	"reflect"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

func fieldName(fld reflect.StructField) string {
	if col, ok := fld.Tag.Lookup("db"); ok && col != "primary" {
		return strings.Split(col, ",")[0]
	}

	if col, ok := fld.Tag.Lookup("json"); ok {
		return strings.Split(col, ",")[0]
	}

	return fld.Name
}

func skipField(i any) bool {
	switch v := i.(type) {
	case IsZeroer:
		if v.IsZero() {
			return true
		}
	case pgtype.Text:
		if !v.Valid {
			return true
		}
	case pgtype.Array[pgtype.Text]:
		if !v.Valid {
			return true
		}
	case pgtype.Timestamptz:
		if !v.Valid {
			return true
		}
	case driver.Valuer:
		if val, err := v.Value(); val == nil || err != nil {
			return true
		}
	}

	return false
}

func writeInt[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](b ByteStringWriter, v T) {
	b.Write(strconv.AppendInt([]byte{}, int64(v), 10))
}

func writeParam(b ByteStringWriter, args *[]any, value any) {
	if col, ok := value.(Columnar); ok {
		col.encodeColumnIdentifier(b)
	} else {
		*args = append(*args, value)
		b.WriteByte('$')
		writeInt(b, len(*args))
	}
}

func writeIdentifier(b ByteStringWriter, identifiers ...string) {
	if len(identifiers) == 0 {
		return
	}

	first := true

	for _, id := range identifiers {
		if first {
			first = false
		} else {
			b.WriteByte('.')
		}

		b.WriteByte('"')
		b.WriteString(id)
		b.WriteByte('"')
	}
}

var stringReplacer = strings.NewReplacer("'", "", "\\", "")

func writeString(b ByteStringWriter, str string) {
	b.WriteByte('\'')
	stringReplacer.WriteString(b, str)
	b.WriteByte('\'')
}
