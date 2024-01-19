package sql

// Any scan any value of a field based on type.
func Any[T, V any](scan func(*T, V)) FieldScanner[T, V] {
	return func(t *T, v V) error {
		scan(t, v)
		return nil
	}
}

// NotNull scan value of a field or apply default value.
func NotNull[T, V any](def V, scan func(*T, V)) FieldScanner[T, *V] {
	return func(t *T, v *V) error {
		if v != nil {
			scan(t, *v)
		} else {
			scan(t, def)
		}

		return nil
	}
}
