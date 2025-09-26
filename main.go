package main

import (
	router "github.com/qww83728/gsam_demo/interface/router"

	"github.com/gin-gonic/gin"
)

var balance = 1000

func main() {
	r := gin.Default()

	router.Router(r)

	r.Run(":8080")
}
