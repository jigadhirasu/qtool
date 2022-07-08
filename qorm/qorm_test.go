package qorm_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestX(t *testing.T) {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MARIADB_USER"),
		os.Getenv("MARIADB_PASS"),
		os.Getenv("MARIADB_HOST"),
		"c_asgame",
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT Doc FROM xxx Limit 10`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var b []byte
		rows.Scan(&b)
		fmt.Println(string(b))
	}

}
