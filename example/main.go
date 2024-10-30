// package main for examples
package main

import (
	"fmt"

	"sqlstream"
	"sqlstream/sql"

	"github.com/jmoiron/sqlx"
)

// Student demo data.
type Student struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Age       int    `db:"age"`
}

func main() {
	db, err := sqlx.Connect("postgres", "user=test dbname=test sslmode=disable")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// create mapping to avoid reflection (db annotations in Student struct just to depict DB schema)
	mapping := sql.StructMapping[Student]{
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

	// make any SQL query
	rows, err := db.Queryx("select * from students")
	if err != nil {
		panic(err)
	}

	// bind mappings with results
	reader, err := sql.NewReader(rows, mapping)
	if err != nil {
		panic(err)
	}

	// read as stream from channel (close on reader and rows done automatically after result is done)
	for student := range sqlstream.Read(reader) {
		fmt.Printf("student: %+v\n", student)
	}
}
