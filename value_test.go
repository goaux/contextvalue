package contextvalue_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/goaux/contextvalue"
)

type A struct {
	Value string
}

func Example() {
	type ID int
	const (
		Red ID = iota
		Blue
	)
	ctx := context.TODO()

	// Set values in the context
	ctx = contextvalue.With(ctx, 42)               // No name
	ctx = contextvalue.With(ctx, &A{Value: "999"}) // No name
	ctx = contextvalue.WithName(ctx, Red, 11)      // With name 'Red'
	ctx = contextvalue.WithName(ctx, Red, "RED")   // Overwrite with string type

	// Retrieve values from the context
	fmt.Println(contextvalue.From[int](ctx))             // Retrieves int value: 42
	fmt.Println(contextvalue.From[*A](ctx))              // Retrieves *A struct: &{999}
	fmt.Println(contextvalue.FromName[int](ctx, Red))    // Retrieves int with name 'Red': 11
	fmt.Println(contextvalue.FromName[string](ctx, Red)) // Retrieves string with name 'Red': "RED"
	fmt.Println(contextvalue.FromName[int](ctx, Blue))   // Retrieves int with name 'Blue': (zero-value and false)

	// // Remove a named value from the context
	ctx = contextvalue.WithoutName[int](ctx, Red)
	fmt.Println(contextvalue.FromName[int](ctx, Red)) // Retrieves int with name 'Red': (zero-value and false)
	// Output:
	// 42 true
	// &{999} true
	// 11 true
	// RED true
	// 0 false
	// 0 false
}

func ExampleWith() {
	ctx := context.TODO()
	ctx = contextvalue.With(ctx, 42)
	value, ok := contextvalue.From[int](ctx)
	fmt.Println(value, ok)
	// Output:
	// 42 true
}

func ExampleWithName() {
	const (
		Red  = "red"
		Blue = "blue"
	)
	ctx := context.TODO()
	ctx = contextvalue.WithName(ctx, Red, 42)
	red, ok := contextvalue.FromName[int](ctx, Red)
	fmt.Println(red, ok)
	blue, ok := contextvalue.FromName[int](ctx, Blue)
	fmt.Println(blue, ok)
	// Output:
	// 42 true
	// 0 false
}

func ExampleWithout() {
	ctx := context.TODO()
	ctx = contextvalue.With(ctx, 42)
	ctx = contextvalue.Without[int](ctx)
	value, ok := contextvalue.From[int](ctx)
	fmt.Println(value, ok)
	// Output:
	// 0 false
}

func ExampleWithoutName() {
	ctx := context.TODO()
	ctx = contextvalue.WithName(ctx, "NAME", 42)
	ctx = contextvalue.WithoutName[int](ctx, "NAME")
	value, ok := contextvalue.FromName[int](ctx, "NAME")
	fmt.Println(value, ok)
	// Output:
	// 0 false
}

func Test(t *testing.T) {
	t.Run("no name", func(t *testing.T) {
		t.Run("int", func(t *testing.T) {
			ctx := context.TODO()
			got, ok := contextvalue.From[int](ctx)
			assertFalse(t, ok, "because the value is not set yet")
			assertZero(t, got, "because ok is false")

			ctx = contextvalue.With(ctx, 42)
			got, ok = contextvalue.From[int](ctx)
			assertTrue(t, ok, "because the value is set")
			assertWant(t, 42, got, "because the value 42 is set")

			got32, ok := contextvalue.From[int32](ctx)
			assertFalse(t, ok, "because the value is not set for type of int32")
			assertZero(t, got32, "because ok is false")

			ctx = contextvalue.With(ctx, 52)
			got, ok = contextvalue.From[int](ctx)
			assertTrue(t, ok, "because the value is set")
			assertWant(t, 52, got, "because the value 52 is set")

			ctx = contextvalue.Without[int](ctx)
			got, ok = contextvalue.From[int](ctx)
			assertFalse(t, ok, "because the value is hidden by Without")
			assertZero(t, got, "because ok is false")
		})
	})

	t.Run("named", func(t *testing.T) {
		type Name int
		const (
			Red Name = iota
			Blue
			Green
		)
		t.Run("int", func(t *testing.T) {
			ctx := context.TODO()
			got, ok := contextvalue.FromName[int](ctx, Red)
			assertFalse(t, ok, "because the value is not set yet")
			assertZero(t, got, "because ok is false")

			ctx = contextvalue.WithName(ctx, Red, 42)
			ctx = contextvalue.WithName(ctx, Blue, 99)

			red, ok := contextvalue.FromName[int](ctx, Red)
			assertTrue(t, ok, "because the value is set")
			assertWant(t, 42, red, "because the value 42 is set")

			blue, ok := contextvalue.FromName[int](ctx, Blue)
			assertTrue(t, ok, "because the value is set")
			assertWant(t, 99, blue, "because the value 99 is set")

			green, ok := contextvalue.FromName[int](ctx, Green)
			assertFalse(t, ok, "because the value is not set for Green")
			assertZero(t, green, "because ok is false")

			ctx = contextvalue.WithoutName[int](ctx, Red)

			red, ok = contextvalue.FromName[int](ctx, Red)
			assertFalse(t, ok, "because the value for Red is hidden by WithoutName")
			assertZero(t, red, "because ok is false")

			blue, ok = contextvalue.FromName[int](ctx, Blue)
			assertTrue(t, ok, "because the value for Blue is not hidden")
			assertWant(t, 99, blue, "because the value 99 is set")
		})
	})
}

func assertTrue(t *testing.T, ok bool, msg ...any) {
	t.Helper()
	if !ok {
		t.Error("must be true", fmt.Sprint(msg...))
	}
}

func assertFalse(t *testing.T, ok bool, msg ...any) {
	t.Helper()
	if ok {
		t.Error("must be false", fmt.Sprint(msg...))
	}
}

func assertZero[T comparable](t *testing.T, v T, msg ...any) {
	t.Helper()
	var zero T
	if v != zero {
		t.Error(fmt.Sprintf("must be zero, got=%v", v), fmt.Sprint(msg...))
	}
}

func assertWant[T comparable](t *testing.T, want, got T, msg ...any) {
	t.Helper()
	if got != want {
		t.Error(fmt.Sprintf("got=%v, want=%v", got, want), fmt.Sprint(msg...))
	}
}
