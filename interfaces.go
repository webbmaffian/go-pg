package pg

import (
	"context"
	"strings"
)

type MutationType uint8

const (
	Inserting MutationType = 1
	Updating  MutationType = 2
)

type BeforeMutation interface {
	BeforeMutation(ctx context.Context, mutationType MutationType) error
}

type AfterMutation interface {
	AfterMutation(ctx context.Context, mutationType MutationType)
}

type IsZeroer interface {
	IsZero() bool
}

type Columnar interface {
	encodeColumn(b *strings.Builder)
	has(col string) bool
}

type AliasedColumnar interface {
	Columnar
	Alias(string) AliasedColumnar
	encodeColumnIdentifier(b *strings.Builder)
}

type MultiColumnar interface {
	Columnar
	Append(...Columnar) MultiColumnar
}

type Condition interface {
	encodeCondition(b *strings.Builder, args *[]any)
}
