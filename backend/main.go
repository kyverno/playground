package main

import (
	"flag"
	"fmt"
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
	var host = flag.String("host", "localhost", "server host")
	var port = flag.Int("port", 8080, "server port")
	var frontend = flag.String("frontend", "../frontend/dist", "frontend folder")

	router := gin.Default()
	router.POST("/engine", engine)
	router.StaticFS("/", http.Dir(*frontend))
	address := fmt.Sprintf("%v:%v", *host, *port)
	if err := router.Run(address); err != nil {
		panic(err)
	}
}
