# contextvalue

`contextvalue` is a Go package that provides type-safe context values using Go generics.
It allows you to store and retrieve strongly typed values in a `context.Context` without using custom key types.

**Important Considerations:**

> Use context Values only for request-scoped data that transits processes and APIs, not for passing optional parameters to functions.

As outlined in the official documentation for [context.WithValue](https://pkg.go.dev/context#WithValue),
context values should **only** be used for request-scoped data that needs to be passed across processes and API boundaries.
Avoid using them for optional function parameters.

## Features

- Type-safe context values using Go generics (requires Go 1.18+)
- No need to define custom key types
- Support for named values with any comparable type as the name
- Helper functions to hide values from the context

## Installation

```sh
go get github.com/goaux/contextvalue
```

## Usage

### Basic Usage

```go
import (
	"context"
	"fmt"
	"github.com/goaux/contextvalue"
)

func main() {
	ctx := context.Background()

	// Store a value in the context
	ctx = contextvalue.With(ctx, 42)

	// Retrieve the value
	value, ok := contextvalue.From[int](ctx)
	fmt.Println(value, ok) // Output: 42 true

	// Type safety: different types have different keys
	strValue, ok := contextvalue.From[string](ctx)
	fmt.Println(strValue, ok) // Output: "" false

	// Hide a value from the context
	ctx = contextvalue.Without[int](ctx)
	value, ok = contextvalue.From[int](ctx)
	fmt.Println(value, ok) // Output: 0 false
}
```

### Named Values

```go
import (
	"context"
	"fmt"
	"github.com/goaux/contextvalue"
)

func main() {
	// Define name types
	type ColorID int
	const (
		Red ColorID = iota
		Blue
		Green
	)

	ctx := context.Background()

	// Store named values
	ctx = contextvalue.WithName(ctx, Red, 42)
	ctx = contextvalue.WithName(ctx, Blue, 99)

	// Retrieve named values
	redValue, ok := contextvalue.FromName[int](ctx, Red)
	fmt.Println(redValue, ok) // Output: 42 true

	blueValue, ok := contextvalue.FromName[int](ctx, Blue)
	fmt.Println(blueValue, ok) // Output: 99 true

	greenValue, ok := contextvalue.FromName[int](ctx, Green)
	fmt.Println(greenValue, ok) // Output: 0 false

	// Hide a named value
	ctx = contextvalue.WithoutName[int](ctx, Red)
	redValue, ok = contextvalue.FromName[int](ctx, Red)
	fmt.Println(redValue, ok) // Output: 0 false
}
```

### Using with Struct Types

```go
import (
	"context"
	"fmt"
	"github.com/goaux/contextvalue"
)

type User struct {
	ID   int
	Name string
}

func main() {
	ctx := context.Background()

	// Store a struct in the context
	ctx = contextvalue.With(ctx, &User{ID: 1, Name: "Alice"})

	// Retrieve the struct
	user, ok := contextvalue.From[*User](ctx)
	fmt.Println(user.Name, ok) // Output: Alice true
}
```

## Why use this package?

The standard Go `context.WithValue` and `context.Value` functions are not type-safe and require you to define custom key types.
This package provides a type-safe API using Go generics, which eliminates the need for custom key types and ensures that you can only retrieve values of the correct type.
This package may not cover all situations, but it should be able to handle most use cases.
