package postgres

import (
	"database/sql"
	"fmt"
	"github.com/EdmundMartin/greenstalk/protocol"
	"github.com/EdmundMartin/greenstalk/stateManager"
	"github.com/lib/pq"
	"log"
	"time"
)

type PGConn struct {
	Db *sql.DB
	Updates chan <- *stateManager.HeapValue
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
	
	CREATE INDEX IF NOT EXISTS tube_state
	ON jobs (tube, state);
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
	return &PGConn{Db: db}
}

func (db *PGConn) Save(j *protocol.Job) (int, error) {
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

func (db *PGConn) Reserve(tubes []string) (*protocol.Job, error) {
	// TODO - IMPLEMENT TTR STATUS RESET
	stmt := `UPDATE jobs SET state = 'RESERVED', until = NOW() + ttr * interval '1 sec' WHERE
            id = (SELECT id FROM jobs WHERE state = 'READY' AND tube = ANY($1) ORDER BY priority ASC LIMIT 1)
			RETURNING id, total_bytes, body, until;`
	for {
		var id, totalBytes int
		var body string
		var until time.Time
		err := db.Db.QueryRow(stmt, pq.Array(tubes)).Scan(&id, &totalBytes, &body, &until)
		if err == nil {
			j := &protocol.Job{ID: id, TotalBytes: totalBytes, Body: body}
			db.Updates <- &stateManager.HeapValue{JobID: id, UnixStamp: until.Unix(), Status: "RESERVED"}
			return j, nil
		} else {
			// TODO break on critical errors
			fmt.Println(err)
		}
		<-time.After(time.Second * 5)
	}
}

func (db *PGConn) Delete(j *protocol.Job) bool {
	var foundID int
	stmt := `DELETE from jobs WHERE id = $1 RETURNING id;`
	err := db.Db.QueryRow(stmt, j.ID).Scan(&foundID)
	if err != nil {
		fmt.Println(err)
	}
	db.Updates <- &stateManager.HeapValue{j.ID, time.Now().Unix(), "DELETED"}
	return foundID == j.ID
}

func (db *PGConn) Bury(j *protocol.Job) bool {
	return false
}

func (db *PGConn) UpdateJob(jID int, status string) {
	switch status {
	case "RESERVED":
		db.resetReserved(jID)
	}
}

func (db *PGConn) resetReserved(jID int) {
	stmt := `UPDATE jobs SET state = 'READY' WHERE id = $1 and state = 'RESERVED';`
	db.Db.Exec(stmt)
}

func (db *PGConn) Reset() {
	stmts := []string{
		`UPDATE jobs SET state = 'READY' WHERE until < NOW() - interval '60 sec' AND state = 'RESERVED';`,
		`UPDATE jobs SET state = 'READY' WHERE until < NOW() - interval '60 sec' AND state = 'DELAYED';`,
	}
	for _, stmt := range stmts {
		db.Db.Exec(stmt)
	}
}