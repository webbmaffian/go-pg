package pg

import (
	"context"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Update(ctx context.Context, db *pgxpool.Pool, table TableSource, src any, condition Condition) (err error) {
	var b strings.Builder
	b.Grow(100)
	args := make([]any, 0, 10)

	b.WriteString("UPDATE ")
	table.encodeQuery(&b, nil)
	b.WriteString(" SET ")

	switch v := src.(type) {

	case map[string]any:
		err = updateFromMap(&b, v, &args)

	case *map[string]any:
		err = updateFromMap(&b, *v, &args)

	case BeforeMutation:
		err = v.BeforeMutation(ctx, Updating)

		if err != nil {
			return
		}

		err = updateFromStruct(&b, src, &args)

	case *BeforeMutation:
		err = (*v).BeforeMutation(ctx, Updating)

		if err != nil {
			return
		}

		err = updateFromStruct(&b, src, &args)

	default:
		err = updateFromStruct(&b, src, &args)
	}

	if err != nil {
		return
	}

	b.WriteByte('\n')
	b.WriteString("WHERE ")
	condition.encodeCondition(&b, &args)

	_, err = db.Exec(ctx, b.String(), args...)

	if err == nil {
		if s, ok := src.(AfterMutation); ok {
			s.AfterMutation(ctx, Updating)
		} else if s, ok := src.(*AfterMutation); ok {
			(*s).AfterMutation(ctx, Updating)
		}
	} else {
		err = QueryError{
			err:   err.Error(),
			query: b.String(),
			args:  args,
		}
	}

	return
}

func updateFromMap(b *strings.Builder, src map[string]any, args *[]any) (err error) {
	first := true

	for k, v := range src {
		if first {
			first = false
		} else {
			b.WriteString(", ")
		}

		writeIdentifier(b, k)
		b.WriteString(" = ")
		writeParam(b, args, v)
	}

	return
}

func updateFromStruct(b *strings.Builder, src any, args *[]any) (err error) {
	first := true
	elem := reflect.ValueOf(src)

	if elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
	}

	typ := elem.Type()
	numFields := elem.NumField()

	for idx := 0; idx < numFields; idx++ {
		f := elem.Field(idx)

		if !f.CanInterface() {
			continue
		}

		fld := typ.Field(idx)
		col := fieldName(fld)

		if fld.Tag.Get("db") == "primary" || col == "-" {
			continue
		}

		v := f.Interface()

		if skipField(v) {
			continue
		}

		if first {
			first = false
		} else {
			b.WriteString(", ")
		}

		writeIdentifier(b, col)
		b.WriteString(" = ")
		writeParam(b, args, v)
	}

	return
}
