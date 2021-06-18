// database provide common database methods
package database

import (
	"fmt"
	"os"
)

func ExampleNewDatabase() {
	os.Setenv("DOCKER_NAME","utest01")
	_, err := NewDatabase()
	fmt.Println(err)

	os.Setenv("DOCKER_NAME","")
	_, err = NewDatabase()
	fmt.Println(err)
	// Output:
	// <nil>
	// <nil>
}

func ExampleExec() {
	db, err := NewDatabase()
	fmt.Println(err)

	_, err = db.Exec("create table if not exists test1 (col1 string)")
	fmt.Println(err)

	_, err = db.Exec("create table error_here")
	fmt.Println(err != nil)
	// Output:
	// <nil>
	// <nil>
	// true
}

func ExampleQuery() {
	db, err := NewDatabase()
	fmt.Println(err)

	_, err = db.Query("select ? + ?", 1, 2)
	fmt.Println(err)

	_, err = db.Query("select error_here", 1, 0)
	fmt.Println(err != nil)
	// Output:
	// <nil>
	// <nil>
	// true
}

