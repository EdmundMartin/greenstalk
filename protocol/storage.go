package protocol

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq" // Required postgres driver
	"log"
)

type Storage interface {
	Save(job *Job) (int, error)
	Delete(job *Job) bool
	Bury(job *Job) bool
	Reserve(tubes []string) (*Job, error)
}


type PGConn struct {
	Db *sql.DB
}

func CreateTable(conn *PGConn) {
	schema := `CREATE TABLE IF NOT EXISTS jobs(
               id SERIAL PRIMARY KEY,
               body TEXT NOT NULL,
               tube varchar(256) NOT NULL,
               priority BIGINT NOT NULL,
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
			 ($1, $2, $3, $4, $5, $6) RETURNING id;`
	err := db.Db.QueryRow(stmt, j.Body, j.Tube, j.Priority, "READY", j.TimeToRun, j.TotalBytes).Scan(&id)
	if err != nil {
		log.Printf("Error inseting new job, %s", err.Error())
		return 0, err
	}
	return id, nil
}

func (db *PGConn) Reserve(tubes []string) (*Job, error) {
	// TODO - IMPLEMENT TTR STATUS RESET
	stmt := `UPDATE jobs SET state = 'RESERVED', until = NOW() + ttr * interval '1 sec' WHERE
            id = (SELECT id FROM jobs WHERE state = 'READY' AND tube = ANY($1) ORDER BY priority ASC LIMIT 1)
			RETURNING id, total_bytes, body;`
	for {
		var id, totalBytes int
		var body string
		err := db.Db.QueryRow(stmt, pq.Array(tubes)).Scan(&id, &totalBytes, &body)
		if err == nil {
			j := &Job{ID:id, TotalBytes:totalBytes, Body:body}
			return j, nil
		} else {
			fmt.Println(err)
		}
	}
}

func (db *PGConn) Delete(j *Job) bool {
	return false
}

func (db *PGConn) Bury(j *Job) bool {
	return false
}
