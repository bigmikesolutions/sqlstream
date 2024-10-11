package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sqlstream/sql"
	"sqlstream/test/db"
)

func Test_ShouldStreamDataFromPostgres(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	pgConn, err := db.ConnectToPostgres(ctx, postgresContainer)
	require.NoError(t, err, "connection error")
	defer func() {
		_ = pgConn.Close()
	}()

	rows, err := pgConn.Queryx(db.SelectAllFromStudents)
	require.NoError(t, err, "select error")

	reader, err := sql.NewReader(rows, db.Mapping)
	require.NoError(t, err, "reader error")

	students := db.ReadAll(reader)

	assert.Equal(t, 3, len(students), "unexpected number of students")
	db.AssertStudent(t, students, db.Student{
		ID:        "1",
		FirstName: "johny",
		LastName:  "bravo",
		Age:       30,
	})
	db.AssertStudent(t, students, db.Student{
		ID:        "2",
		FirstName: "mike",
		LastName:  "tyson",
		Age:       51,
	})
	db.AssertStudent(t, students, db.Student{
		ID:        "3",
		FirstName: "pamela",
		LastName:  "anderson",
		Age:       65,
	})
}
