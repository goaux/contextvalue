// Package contextvalue provides type-safe context values using Go generics.
// It allows storing and retrieving strongly typed values in a context.Context.
package contextvalue

import "context"

type key[T any] struct{}

// With stores a value of type T in the provided context and returns the updated context.
// It uses a zero-value struct parameterized with the type T as the context key.
func With[T any](ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, key[T]{}, value)
}

// From retrieves a value of type T from the provided context.
// It returns the value and a boolean indicating whether the value was found.
func From[T any](ctx context.Context) (T, bool) {
	v, ok := ctx.Value(key[T]{}).(T)
	return v, ok
}

type keyName[T any, N comparable] struct {
	name N
}

// WithName stores a named value of type T in the provided context and returns the updated context.
// It uses a struct containing the name of type N as the context key.
// The name type N must be comparable.
func WithName[T any, N comparable](ctx context.Context, name N, value T) context.Context {
	return context.WithValue(ctx, keyName[T, N]{name: name}, value)
}

// FromName retrieves a named value of type T from the provided context.
// It returns the value and a boolean indicating whether the value was found.
func FromName[T any, N comparable](ctx context.Context, name N) (T, bool) {
	v, ok := ctx.Value(keyName[T, N]{name: name}).(T)
	return v, ok
}

// Without hides a value of type T from the provided context by setting it to nil.
// It returns the updated context.
func Without[T any](ctx context.Context) context.Context {
	return context.WithValue(ctx, key[T]{}, nil)
}

// WithoutName hides a named value of type T from the provided context by setting it to nil.
// It returns the updated context.
func WithoutName[T any, N comparable](ctx context.Context, name N) context.Context {
	return context.WithValue(ctx, keyName[T, N]{name: name}, nil)
}
