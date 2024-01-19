package sql

type (
	// FieldSetter is a function that sets a value in a target.
	FieldSetter[T any] func(*T) error

	// Scanner column scanner.
	// Similar to https://pkg.go.dev/database/sql#Scanner but wit source and mapper move outisde.
	Scanner[T any] interface {
		Scan() (any, FieldSetter[T])
	}

	// FieldScanner defines a scanner rule for a column.
	FieldScanner[T, V any] func(*T, V) error

	// StructMapping defines rules for mapping columns to a specific type.
	StructMapping[T any] map[string]Scanner[T]
)

// Scan implements the Scan method a chosen field of a given type.
func (f FieldScanner[T, V]) Scan() (any, FieldSetter[T]) {
	var v V
	return &v, func(t *T) error {
		return f(t, v)
	}
}
