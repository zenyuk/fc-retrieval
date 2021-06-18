// database provide common database methods
package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//
type Database struct {
	db *sql.DB
}

// NewDatabase returns
func NewDatabase() (*Database, error) {
	_ = os.MkdirAll("logs", os.ModePerm)

	host, _ := os.Hostname()
	dockerName := os.Getenv("DOCKER_NAME")

	dbFileName := ".sqlite"
	if dockerName == "" {
		dbFileName = host + dbFileName
	} else {
		dbFileName = dockerName + dbFileName
	}
	db, err := sql.Open("sqlite3", "logs" + "/" + dbFileName)
	return &Database{db:db}, err
}

//
func (db *Database) Exec (statement string, parameters ...interface{}) (res sql.Result, err error) {
	defer func() {
		// fmt.Printf("`%v`\t`%v`\n", err, statement)
	}()
	stmt, err := db.db.Prepare(statement)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err = stmt.Exec(parameters...)
	return res, err
}

func (db *Database) Query (statement string, parameters ...interface{}) (res *sql.Rows, err error) {
	defer func() {
		// fmt.Printf("`%v`\t`%v`\n", err, statement)
	}()
	stmt, err := db.db.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	res, err = stmt.Query(parameters...)
	return res, err
}

