package test

import (
	"context"
	"testing"
	"time"

	"sqlstream/sql"
	"sqlstream/test/containers/pg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ShouldStreamDataFromPostgres(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	pgConn, err := dc.NewDB(ctx)
	require.NoError(t, err, "connection error")
	defer func() {
		_ = pgConn.Close()
	}()

	rows, err := pgConn.Queryx(pg.SelectAllFromStudents)
	require.NoError(t, err, "select error")

	reader, err := sql.NewReader(rows, pg.Mapping)
	require.NoError(t, err, "reader error")

	students := pg.ReadAll(reader)

	assert.Equal(t, 3, len(students), "unexpected number of students")
	pg.AssertStudent(t, students, pg.Student{
		ID:        "1",
		FirstName: "johny",
		LastName:  "bravo",
		Age:       30,
	})
	pg.AssertStudent(t, students, pg.Student{
		ID:        "2",
		FirstName: "mike",
		LastName:  "tyson",
		Age:       51,
	})
	pg.AssertStudent(t, students, pg.Student{
		ID:        "3",
		FirstName: "pamela",
		LastName:  "anderson",
		Age:       65,
	})
}

func Test_ShouldStreamDataFromPostgresWithContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	pgProxy, err := dc.NewDBProxy(ctx)
	require.NoError(t, err, "connection error")
	defer pgProxy.Close()

	ctx, cancel = context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	rows, err := pgProxy.DB.QueryxContext(ctx, pg.SelectAllFromStudents)
	require.NoError(t, err, "select error")

	reader, err := sql.NewReader(rows, pg.Mapping)
	require.NoError(t, err, "reader error")

	students := pg.ReadAll(reader)

	assert.Equal(t, 3, len(students), "unexpected number of students")
	pg.AssertStudent(t, students, pg.Student{
		ID:        "1",
		FirstName: "johny",
		LastName:  "bravo",
		Age:       30,
	})
	pg.AssertStudent(t, students, pg.Student{
		ID:        "2",
		FirstName: "mike",
		LastName:  "tyson",
		Age:       51,
	})
	pg.AssertStudent(t, students, pg.Student{
		ID:        "3",
		FirstName: "pamela",
		LastName:  "anderson",
		Age:       65,
	})
}
