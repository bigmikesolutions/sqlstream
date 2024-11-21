package sql_test

import "sqlstream/sql"

type (
	testStruct struct {
		String  string
		StringP *string
		Int     int
		Float32 float32
	}
)

var testStructColumns = []string{"string", "stringp", "int", "float32"}

var testStructMapping = sql.StructMapping[testStruct]{
	"string":  sql.Any(func(t *testStruct, v string) { t.String = v }),
	"stringp": sql.Any(func(t *testStruct, v *string) { t.StringP = v }),
	"int":     sql.Any(func(t *testStruct, v int) { t.Int = v }),
	"float32": sql.Any(func(t *testStruct, v float32) { t.Float32 = v }),
}

func newTestStructRows(data [][]any) *mockRows {
	return &mockRows{
		idx:     -1,
		data:    data,
		columns: testStructColumns,
	}
}

func ptr[T any](v T) *T {
	return &v
}
