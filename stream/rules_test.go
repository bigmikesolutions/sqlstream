package stream_test

import (
	"testing"

	"github.com/bigmikesolutions/sqlstream/stream"
	"github.com/stretchr/testify/assert"
)

func Test_Rules(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name           string
		rows           stream.Rows
		columns        stream.StructMapping[testStruct]
		expRows        []testStruct
		expReaderErr   error
		expScanningErr error
	}

	tests := []testCase{
		{
			name: "any",
			rows: newTestStructRows([][]any{
				{"v11", ptr("v12"), 1, 1.3},
				{"v21", ptr("v22"), 2, 2.3},
			}),
			columns: stream.StructMapping[testStruct]{
				"string":  stream.Any(func(t *testStruct, v string) { t.String = v }),
				"stringp": stream.Any(func(t *testStruct, v *string) { t.StringP = v }),
				"int":     stream.Any(func(t *testStruct, v int) { t.Int = v }),
				"float32": stream.Any(func(t *testStruct, v float32) { t.Float32 = v }),
			},
			expRows: []testStruct{
				{"v11", ptr("v12"), 1, 1.3},
				{"v21", ptr("v22"), 2, 2.3},
			},
		},
		{
			name: "not null",
			rows: newTestStructRows([][]any{
				{"v00", nil, nil, nil},
				{"v11", ptr("v12"), ptr(1), ptr(1.3)},
				{"v21", "v22", ptr(2), ptr(2.3)},
			}),
			columns: stream.StructMapping[testStruct]{
				"string":  stream.NotNull("def1", func(t *testStruct, v string) { t.String = v }),
				"stringp": stream.NotNull("def2", func(t *testStruct, v string) { t.StringP = &v }),
				"int":     stream.NotNull(-1, func(t *testStruct, v int) { t.Int = v }),
				"float32": stream.NotNull(-2.5, func(t *testStruct, v float32) { t.Float32 = v }),
			},
			expRows: []testStruct{
				{"v00", ptr("def2"), -1, -2.5},
				{"v11", ptr("v12"), 1, 1.3},
				{"v21", ptr("v22"), 2, 2.3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reader, err := stream.NewReader(tt.rows, tt.columns)
			assert.Equalf(t, tt.expReaderErr, err, "unexpected new reader err: %s", err)

			if reader != nil {
				defer func() {
					_ = reader.Close()
				}()

				var results []testStruct

				for reader.Next() {
					var o testStruct
					err := reader.Scan(&o)
					assert.Equalf(t, tt.expScanningErr, err, "unexpected scanning err: %s", err)
					if err == nil {
						results = append(results, o)
					}
				}

				assert.Equal(t, tt.expRows, results, "unexpected rows")
			}
		})
	}
}
