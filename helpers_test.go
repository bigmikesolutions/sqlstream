package sqlstream_test

import "sync"

type mockTRows[T any] struct {
	mu sync.Mutex

	idx     int
	data    []T
	err     error
	scanErr error

	closed   bool
	closeErr error
}

func newMockTRows[T any](data []T) *mockTRows[T] {
	return &mockTRows[T]{
		idx:  -1,
		data: data,
	}
}

func (m *mockTRows[T]) Next() bool {
	m.idx++
	return m.idx < len(m.data)
}

func (m *mockTRows[T]) Scan(t *T) error {
	if m.scanErr != nil {
		return m.scanErr
	}

	*t = m.data[m.idx]
	return nil
}

func (m *mockTRows[T]) Err() error {
	return m.err
}

func (m *mockTRows[T]) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return m.closeErr
}

func (m *mockTRows[T]) Closed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closed
}
