package pg

var _ OrderByColumnar = OrderBy{}

type OrderBy []OrderByColumnar

func (o OrderBy) IsZero() bool {
	for i := range o {
		if !o[i].IsZero() {
			return false
		}
	}

	return true
}

func (o OrderBy) encodeOrderBy(b ByteStringWriter) {
	if len(o) == 0 {
		return
	}

	o[0].encodeOrderBy(b)

	for _, v := range o[1:] {
		b.WriteString(", ")
		v.encodeOrderBy(b)
	}
}

func Asc(cols ...any) OrderByColumnar {
	if cols == nil {
		return nil
	}

	columns := make(asc, len(cols))

	for i := range cols {
		switch c := cols[i].(type) {

		case Columnar:
			columns[i] = c

		case string:
			columns[i] = Column(c)
		}
	}

	return columns
}

type asc []Columnar

func (o asc) IsZero() bool {
	for i := range o {
		if !o[i].IsZero() {
			return false
		}
	}

	return true
}

func (o asc) encodeOrderBy(b ByteStringWriter) {
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

func Desc(cols ...any) OrderByColumnar {
	if cols == nil {
		return nil
	}

	columns := make(desc, len(cols))

	for i := range cols {
		switch c := cols[i].(type) {

		case Columnar:
			columns[i] = c

		case string:
			columns[i] = Column(c)
		}
	}

	return columns
}

type desc []Columnar

func (o desc) IsZero() bool {
	for i := range o {
		if !o[i].IsZero() {
			return false
		}
	}

	return true
}

func (o desc) encodeOrderBy(b ByteStringWriter) {
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
