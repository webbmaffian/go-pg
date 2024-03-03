package pg

import (
	"fmt"
	"testing"

	"github.com/webbmaffian/go-fast"
)

func Example_printfToQuery() {
	buf := fast.NewStringBuffer(128)
	args := printfToQuery(buf, "SELECT * WHERE id = %3.f AND foo = %s AND bar = %[3]d AND baz = %d", nil)
	fmt.Println("query:", buf.String())
	fmt.Println("args:", args)

	// Output: Wazzaaaa
}

func Example_printfToQuery_args() {
	buf := fast.NewStringBuffer(128)
	args := printfToQuery(buf, "SELECT * FROM %T WHERE id = %d AND foo = %[1]s AND bar = %d AND baz = %d", []any{Table(nil, "users")})
	fmt.Println("query:", buf.String())
	fmt.Println("args:", args)

	// Output: Wazzaaaa
}

func Benchmark_printfToQuery(b *testing.B) {
	buf := fast.NewStringBuffer(128)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = printfToQuery(buf, "SELECT * WHERE id = %3.f AND foo = %s AND bar = %[3]d AND baz = %d", nil)
		buf.Reset()
	}
}

func Benchmark_printfToQuery_args(b *testing.B) {
	buf := fast.NewStringBuffer(128)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = printfToQuery(buf, "SELECT * FROM %T WHERE id = %3.f AND foo = %s AND bar = %[3]d AND baz = %d", []any{Table(nil, "users")})
		buf.Reset()
	}
}
