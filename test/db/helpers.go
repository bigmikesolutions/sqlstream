package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sqlstream"
	"sqlstream/sql"
)

func ReadAll(reader sql.TRows[Student]) map[string]Student {
	students := make(map[string]Student)
	for student := range sqlstream.Read(reader) {
		students[student.Value.ID] = student.Value
	}
	return students
}

func AssertStudent(t *testing.T, students map[string]Student, exp Student) {
	student, ok := students[exp.ID]
	assert.Truef(t, ok, "student 1 was not found")
	if ok {
		assert.Equalf(t, exp, student, "unexpected student: %+v", exp)
	}
}
