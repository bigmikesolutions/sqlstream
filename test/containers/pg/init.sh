#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL

 CREATE TABLE students (
                       id SERIAL,
                       first_name TEXT NOT NULL,
                       last_name TEXT NOT NULL,
                       age INT NOT NULL
);

INSERT INTO students(id,first_name,last_name,age) values
('1','johny','bravo',30),
('2','mike','tyson',51),
('3','pamela','anderson',65);

EOSQL