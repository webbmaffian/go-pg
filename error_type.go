package pg

import "fmt"

var _ error = Error("foobar")
var _ fmt.Stringer = Error("foobar")

type Error string

func (err Error) Error() string {
	return string(err)
}

func (err Error) String() string {
	return string(err)
}
