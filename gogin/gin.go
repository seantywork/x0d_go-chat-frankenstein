package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthiness struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// albums slice to seed record album data.
var healthiness_record = healthiness{
	ID: "_ptt_healthiness_probe", Status: "Healthy",
}

func main() {
	router := gin.Default()

	router.GET("/getHealth", getHealth)

	router.POST("/postHealth", postHealth)

	router.Run("0.0.0.0:8080")
}

func getHealth(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, healthiness_record)
}

func postHealth(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, healthiness_record)
}
