package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type engineRequest struct {
	Policy   string `json:"policy"`
	Resource string `json:"resource"`
	Context  string `json:"context"`
}

type engineResponse struct {
}

func engine(c *gin.Context) {
	var request engineRequest
	if err := c.BindJSON(&request); err != nil {
		return
	}
	var response engineResponse
	c.IndentedJSON(http.StatusCreated, response)
}

func main() {
	router := gin.Default()
	router.POST("/engine", engine)
	router.StaticFS("/", http.Dir("../frontend/dist"))
	if err := router.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
