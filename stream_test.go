package sqlstream_test

import (
	"errors"
	"testing"

	"sqlstream"

	"github.com/stretchr/testify/assert"
)

type testObject struct {
	ID     int
	Field1 string
}

func TestStream_ShouldReadRows(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name      string
		in        *mockTRows[testObject]
		expResult []sqlstream.Entry[testObject]
	}

	tests := []testCase{
		{
			name: "read rows",
			in: newMockTRows[testObject]([]testObject{
				{1, "1"},
				{2, "2"},
			}),
			expResult: []sqlstream.Entry[testObject]{
				{Value: testObject{1, "1"}},
				{Value: testObject{2, "2"}},
			},
		},
		{
			name: "read rows errors",
			in: &mockTRows[testObject]{
				scanErr: errors.New("scan error"), // nolint:all
				idx:     -1,
				data: []testObject{
					{1, "1"},
					{2, "2"},
				},
			},
			expResult: []sqlstream.Entry[testObject]{
				{Err: errors.New("scan error")}, // nolint:all
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var results []sqlstream.Entry[testObject]

			for entry := range sqlstream.Read[testObject](tt.in) {
				results = append(results, entry)
			}

			assert.Equal(t, tt.expResult, results)
			assert.True(t, tt.in.Closed(), "expected rows to be closed")
		})
	}
}
