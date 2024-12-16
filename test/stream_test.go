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

const rowsInDb = 3 // loaded by init script

func Test_ShouldStreamDataFromPostgres(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	pgConn, err := dc.DB(ctx)
	require.NoError(t, err, "connection error")
	defer func() {
		_ = pgConn.Close()
	}()

	rows, err := pgConn.Queryx(pg.SelectAllFromStudents)
	require.NoError(t, err, "select error")

	reader, err := sql.NewReader(rows, pg.Mapping)
	require.NoError(t, err, "reader error")

	students := pg.ReadAll(reader)

	assert.Equal(t, rowsInDb, len(students), "unexpected number of students")
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

func Test_ShouldNoStreamDataForQueryTimeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	pgProxy, err := dc.DBProxy(ctx)
	require.NoError(t, err, "connection error")
	defer pgProxy.Close()

	err = pgProxy.Upstream.AddLatency(1000, 100)
	require.NoError(t, err, "proxy latency error")

	queryCtx, queryCancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer queryCancel()

	rows, err := pgProxy.DB.QueryxContext(queryCtx, pg.SelectAllFromStudents)
	require.Error(t, err, "select must fail due to timeout caused by latency")
	require.Nil(t, rows, "no rows must be returned for timeouts")
}
