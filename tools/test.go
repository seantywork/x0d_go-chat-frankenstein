package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type user_data struct {
	DATA []user_record
}

type user_record struct {
	UUID    string
	USID    string
	USER_ID string
	USER_PW string
	ACTIVE  int
}

func dbQuery(query string, args []string) (interface{}, error) {

	db, err := sql.Open("mysql", "seantywork:youdonthavetoknow@tcp(127.0.0.1:3306)/chfrank")

	if err != nil {
		return -1, err
	}

	for i := 0; i < len(args); i++ {

		query = strings.Replace(query, "?", args[i], 1)

	}

	defer db.Close()

	results, err := db.Query(query)

	if err != nil {

		return 1, err

	}

	return results, err

}

func main() {

	var result_rows *sql.Rows

	q := "SELECT * FROM chfrank_user WHERE user_id LIKE '%?%'"

	a := []string{"test"}

	res, err := dbQuery(q, a)

	if err != nil {

		panic(err.Error())

	}

	result_rows = res.(*sql.Rows)

	for result_rows.Next() {
		var record user_record

		err = result_rows.Scan(&record.UUID, &record.USID, &record.USER_ID, &record.USER_PW, &record.ACTIVE)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(record.UUID + " " + record.USER_ID + " " + record.USER_PW)

	}

}