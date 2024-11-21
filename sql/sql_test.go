package sql_test

type mockRows struct {
	err error

	closed   bool
	closeErr error

	columns    []string
	columnsErr error

	data    [][]any
	idx     int
	scanErr error
}

func (m *mockRows) Next() bool {
	m.idx++
	return m.idx < len(m.data)
}

func (m *mockRows) Scan(dest ...any) error {
	if m.scanErr != nil {
		return m.scanErr
	}

	for index, d := range dest {
		switch value := d.(type) {

		case *string:
			switch s := m.data[m.idx][index].(type) {
			case string:
				*value = s
			case *string:
				*value = *s
			}

		case *int64:
			switch s := m.data[m.idx][index].(type) {
			case int64:
				*value = s
			case int:
				*value = int64(s)
			}

		case *int:
			switch s := m.data[m.idx][index].(type) {
			case int:
				*value = s
			case int64:
				*value = int(s)
			}

		case *float32:
			switch s := m.data[m.idx][index].(type) {
			case float32:
				*value = s
			case float64:
				*value = float32(s)
			}

		case **string:
			switch s := m.data[m.idx][index].(type) {
			case string:
				*value = &s
			case *string:
				*value = s
			}

		case **int:
			switch s := m.data[m.idx][index].(type) {
			case int:
				*value = &s

			case *int:
				*value = s
			}

		case **float32:
			switch s := m.data[m.idx][index].(type) {
			case *float32:
				*value = s
			case *float64:
				v := float32(*s)
				*value = &v
			}

		}
	}

	return nil
}

func (m *mockRows) Columns() ([]string, error) {
	return m.columns, m.columnsErr
}

func (m *mockRows) Err() error {
	return m.err
}

func (m *mockRows) Close() error {
	m.closed = true
	return m.closeErr
}
