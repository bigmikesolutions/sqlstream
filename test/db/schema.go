package db

import (
	"sqlstream/sql"
)

const (
	SelectAllFromStudents = `SELECT * FROM students`
)

type Student struct {
	ID        string `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Age       int    `db:"age"`
}

var Mapping = sql.StructMapping[Student]{
	"id": sql.Any(func(s *Student, v string) {
		s.ID = v
	}),
	"first_name": sql.Any(func(s *Student, v string) {
		s.FirstName = v
	}),
	"last_name": sql.Any(func(s *Student, v string) {
		s.LastName = v
	}),
	"age": sql.Any(func(s *Student, v int) {
		s.Age = v
	}),
}
