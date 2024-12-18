// Package sqlstream provides stream capability for sql.
package sqlstream

import "github.com/bigmikesolutions/sqlstream/sql"

type (
	// Entry single results from a stream.
	Entry[T any] struct {
		Value T
		Err   error
	}

	// ReadStream read-only stream of results from a query.
	ReadStream[T any] <-chan Entry[T]
)

// Read creates a read-only stream of results from a query.
func Read[T any](in sql.TRows[T]) ReadStream[T] {
	out := make(chan Entry[T])

	go func() {
		defer func() { _ = in.Close() }()
		defer close(out)

		for in.Next() {
			var t T

			if err := in.Scan(&t); err != nil {
				out <- Entry[T]{Err: err}
				break
			}

			out <- Entry[T]{Value: t}
		}
	}()

	return out
}
