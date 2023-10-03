package main

import (
	"context"
	"fmt"
)

func step1(ctx context.Context) context.Context {
	child := context.WithValue(ctx, "name", "大脸猫")
	return child
}
func step2(ctx context.Context) context.Context {
	child := context.WithValue(ctx, "age", 18)
	return child
}
func step3(ctx context.Context) {
	fmt.Printf("n %s\n", ctx.Value("name"))
	fmt.Printf("a %d\n", ctx.Value("age"))
}
func main() {
	ctx := context.Background()
	ctx = step1(ctx)
	ctx = step2(ctx)
	step3(ctx)
}
