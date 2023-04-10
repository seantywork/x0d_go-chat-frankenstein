package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type healthiness struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type user_data struct {
	DATA []user_record
}

type user_record struct {
	UUID    string
	USER_ID string
	USER_PW string
	ACTIVE  int
}

// albums slice to seed record album data.
var healthiness_record = healthiness{
	ID: "_ptt_healthiness_probe", Status: "Healthy",
}

func main() {
	router := gin.Default()

	router.GET("/getHealth", getHealth)

	router.POST("/postHealth", postHealth)

	router.GET("/dbCheck", dbCheck)

	router.Run("0.0.0.0:8888")
}

func getHealth(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, healthiness_record)
}

func postHealth(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, healthiness_record)
}

func dbCheck(c *gin.Context) {

	var records user_data

	db, err := sql.Open("mysql", "seantywork:youdonthavetoknow@tcp(127.0.0.1:3306)/chfrank")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	results, err := db.Query("SELECT * FROM chfrank_user")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var record user_record

		err = results.Scan(&record.UUID, &record.USER_ID, &record.USER_PW, &record.ACTIVE)
		if err != nil {
			panic(err.Error())
		}

		records.DATA = append(records.DATA, record)

	}

	c.IndentedJSON(http.StatusOK, records)

}
