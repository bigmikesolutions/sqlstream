[![codecov](https://codecov.io/gh/bigmikesolutions/sqlstream/graph/badge.svg?token=LCSGUE3GAG)](https://codecov.io/gh/bigmikesolutions/sqlstream)

# sqlstream

Small library to handle SQL queries as streams (channels).

Main concerns:

* Mapping is done in functional fashion without reflection.

* Memory usage: mapping and result appears row by row instead all at once

* One-to-many/many-to-many even for nested structures

## Benchmark results

First implementation required 1 memory allocation per operation (row result coming from DB), i.e.: 

```go
pkg: sqlstream
BenchmarkStream_ReadRows-10     14924950               150.9 ns/op            24 B/op          1 allocs/op
BenchmarkStream_ReadRows-10     15875574               150.4 ns/op            24 B/op          1 allocs/op
```

## Usage

All samples of usage can be found in `example` directory.

But main usage looks like following:

```go
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
```