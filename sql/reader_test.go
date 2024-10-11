package sql_test

import (
	"errors"
	"fmt"
	"testing"

	"sqlstream/sql"

	"github.com/stretchr/testify/assert"
)

// nolint: err113
func TestNewReader_ShouldReadRows(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name           string
		rows           sql.Rows
		columns        sql.StructMapping[testStruct]
		expRows        []testStruct
		expReaderErr   error
		expScanningErr error
	}

	tests := []testCase{
		{
			name: "read rows",
			rows: newTestStructRows([][]any{
				{"v11", stringP("v12"), 1, 1.3},
				{"v21", stringP("v22"), 2, 2.3},
			}),
			columns: testStructMapping,
			expRows: []testStruct{
				{"v11", stringP("v12"), 1, 1.3},
				{"v21", stringP("v22"), 2, 2.3},
			},
		},
		{
			name: "handle error while getting column names",
			rows: &mockRows{
				columnsErr: errors.New("columns error"), // nolint:all
			},
			columns:      testStructMapping,
			expReaderErr: fmt.Errorf("column names: %w", errors.New("columns error")), // nolint:all
		},
		{
			name: "handle error while scanning",
			rows: &mockRows{
				scanErr: errors.New("columns error"), // nolint:all
			},
			columns:        testStructMapping,
			expScanningErr: fmt.Errorf("column names: %w", errors.New("columns error")), // nolint:all
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reader, err := sql.NewReader(tt.rows, tt.columns)
			assert.Equalf(t, tt.expReaderErr, err, "unexpected new reader err: %s", err)

			if reader != nil {
				defer reader.Close()

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
