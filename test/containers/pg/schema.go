package pg

import "github.com/bigmikesolutions/sqlstream/stream"

const (
	// SelectAllFromStudents test query to get all data from students table.
	SelectAllFromStudents = `SELECT * FROM students`
)

// Student student with test data.
type Student struct {
	ID        string `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Age       int    `db:"age"`
}

// Mapping mapping for student data.
var Mapping = stream.StructMapping[Student]{
	"id": stream.Any(func(s *Student, v string) {
		s.ID = v
	}),
	"first_name": stream.Any(func(s *Student, v string) {
		s.FirstName = v
	}),
	"last_name": stream.Any(func(s *Student, v string) {
		s.LastName = v
	}),
	"age": stream.Any(func(s *Student, v int) {
		s.Age = v
	}),
}
