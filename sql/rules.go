package sql

// Any creates a custom scanner for a column with a specified scan function.
func Any[T, V any](scan func(*T, V)) FieldScanner[T, V] {
	return func(t *T, v V) error {
		scan(t, v)
		return nil
	}
}

// Null creates a custom scanner for a nullable column with a specified default value and scan function.
func Null[T, V any](def V, scan func(*T, V)) FieldScanner[T, *V] {
	return func(t *T, v *V) error {
		if v != nil {
			scan(t, *v)
		} else {
			scan(t, def)
		}

		return nil
	}
}
