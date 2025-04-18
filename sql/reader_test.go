package sql_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bigmikesolutions/sqlstream/sql"

	"github.com/stretchr/testify/assert"
)

// nolint: err113
func Test_Reader_ShouldReadRows(t *testing.T) {
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
				{"v11", ptr("v12"), 1, 1.3},
				{"v21", ptr("v22"), 2, 2.3},
			}),
			columns: testStructMapping,
			expRows: []testStruct{
				{"v11", ptr("v12"), 1, 1.3},
				{"v21", ptr("v22"), 2, 2.3},
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reader, err := sql.NewReader(tt.rows, tt.columns)
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
