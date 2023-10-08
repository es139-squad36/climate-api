package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type renewableData struct {
	ID string
}

var data = []renewableData{
	{ID: "1"},
}

func getData(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, data)
}

func addData(context *gin.Context) {
	var entry renewableData

	if err := context.BindJSON(&entry); err != nil {
		return
	}

	data = append(data, entry)
	context.IndentedJSON(http.StatusCreated, entry)
}

func main() {
	router := gin.Default()
	router.GET("/getData", getData)
	router.POST("/addData", addData)
	router.Run()
}
