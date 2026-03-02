package stream_test

import "github.com/bigmikesolutions/sqlstream/stream"

type (
	testStruct struct {
		String  string
		StringP *string
		Int     int
		Float32 float32
	}
)

var testStructColumns = []string{"string", "stringp", "int", "float32"}

var testStructMapping = stream.StructMapping[testStruct]{
	"string":  stream.Any(func(t *testStruct, v string) { t.String = v }),
	"stringp": stream.Any(func(t *testStruct, v *string) { t.StringP = v }),
	"int":     stream.Any(func(t *testStruct, v int) { t.Int = v }),
	"float32": stream.Any(func(t *testStruct, v float32) { t.Float32 = v }),
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
