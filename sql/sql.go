package sql

type (
	// Rows is an interface defined based on sql rows: https://pkg.go.dev/database/sql#Rows to decouple
	// from other libs like sql or sqlx.
	Rows interface {
		// Next prepares the next result row for reading with the Scan method.
		Next() bool
		// Scan copies the columns in the current row into the values pointed at by dest. The number of values in dest must be the same as the number of columns in Rows.
		Scan(dest ...any) error
		// StructMapping returns the column names. StructMapping returns an error if the rows are closed.
		Columns() ([]string, error)
		// Err returns the error, if any, that was encountered during iteration. Err may be called after an explicit or implicit Close.
		Err() error
		// Close closes the Rows, preventing further enumeration.
		Close() error
	}

	// TRows represents sql rows with a specific type.
	TRows[T any] interface {
		// Next prepares the next result row for reading with the Scan method.
		Next() bool
		// Scan copies the columns in the current row of type T.
		Scan(t *T) error
		// Err returns the error, if any, that was encountered during iteration. Err may be called after an explicit or implicit Close.
		Err() error
		// Close closes the TRows, preventing further enumeration.
		Close() error
	}
)
