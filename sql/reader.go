package sql

import (
	"fmt"
)

type reader[T any] struct {
	rows      Rows
	rowValues []any
	setters   []FieldSetter[T]
}

// NewReader creates a new row reader with column mappings for concrete type T.
func NewReader[T any](rows Rows, columns StructMapping[T]) (TRows[T], error) {
	names, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("column names: %w", err)
	}

	rowValues := make([]any, len(names))
	setters := make([]FieldSetter[T], len(names))

	for i, n := range names {
		if s, ok := columns[n]; ok {
			rowValues[i], setters[i] = s.Scan()
		} else {
			rowValues[i] = new(any)
		}
	}

	return reader[T]{
		rows:      rows,
		rowValues: rowValues,
		setters:   setters,
	}, nil
}

func (i reader[T]) Close() error {
	return i.rows.Close()
}

func (i reader[T]) Err() error {
	return i.rows.Err()
}

func (i reader[T]) Next() bool {
	return i.rows.Next()
}

func (i reader[T]) Scan(t *T) error {
	if err := i.rows.Scan(i.rowValues...); err != nil {
		return err
	}

	for _, setter := range i.setters {
		if setter != nil {
			if err := setter(t); err != nil {
				return err
			}
		}
	}

	return nil
}
