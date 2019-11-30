package protocol

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Required postgres driver
	"log"
)

type Storage interface {
	Save(job *Job) (int, error)
	Delete(job *Job) bool
	Bury(job *Job) bool
}


type PGConn struct {
	Db *sql.DB
}

func createTable(conn *PGConn) {
	schema := `CREATE TABLE IF NOT EXISTS jobs(
               id SERIAL PRIMARY KEY,
               body TEXT NOT NULL,
               tube varchar(256) NOT NULL,
               priority INTEGER NOT NULL,
               state varchar(40) NOT NULL,
               ttr INTEGER NOT NULL,
               total_bytes INTEGER,
               until TIMESTAMP
	);
	`
	_, err := conn.Db.Exec(schema)
	if err != nil {
		log.Fatalln(err)
	}
}

func NewPGConn(host, user, pwd, dbname string, port int) *PGConn {
	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pwd, dbname)
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Fatal(err)
	}
	return &PGConn{Db:db}
}

func (db *PGConn) Save(j *Job) (int, error) {
	var id int
	stmt := `INSERT INTO jobs (body, tube, priority, state, ttr, total_bytes) VALUES
			 ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	err := db.Db.QueryRow(stmt, j.Body, j.Tube, j.Priority, "READY", j.TimeToRun).Scan(&id)
	if err != nil {
		log.Printf("Error inseting new job, %s", err.Error())
		return 0, err
	}
	return id, nil
}

func (db *PGConn) Delete(j *Job) bool {
	return false
}

func (db *PGConn) Bury(j *Job) bool {
	return false
}
