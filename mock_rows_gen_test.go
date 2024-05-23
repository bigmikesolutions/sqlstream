package sqlstream_test

type rowGenerator[T any] func() (*T, error)

type rowsGenerator[T any] struct {
	gen     rowGenerator[T]
	next    *T
	nextErr error
}

func newTRowsGenerator[T any](gen rowGenerator[T]) *rowsGenerator[T] {
	return &rowsGenerator[T]{
		gen: gen,
	}
}

func (m *rowsGenerator[T]) Next() bool {
	m.next, m.nextErr = m.gen()
	return m.nextErr == nil && m.next != nil
}

func (m *rowsGenerator[T]) Scan(t *T) error {
	*t = *m.next
	return nil
}

func (m *rowsGenerator[T]) Err() error {
	return m.nextErr
}

func (m *rowsGenerator[T]) Close() error {
	return nil
}

func (m *rowsGenerator[T]) Closed() bool {
	return true
}
