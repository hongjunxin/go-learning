package main

import (
	"context"
	"fmt"
)

type Span struct {
	number int
}

type Transaction struct {
	name string
}

type spanKey struct{}
type transactionKey struct{}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, spanKey{}, Span{1}) // set kv to ctx
	ctx = context.WithValue(ctx, transactionKey{}, Transaction{"transaction"})
	ctx = context.WithValue(ctx, spanKey{}, Span{2})

	// kv 原则，相同的 key 会被最后设置的覆盖掉
	fmt.Printf("span.number=%v\n", ctx.Value(spanKey{}).(Span).number) // output 2
	fmt.Printf("span.number=%v\n", ctx.Value(spanKey{}).(Span).number) // output 2
	fmt.Printf("transaction.name=%v\n", ctx.Value(transactionKey{}).(Transaction).name)
}
