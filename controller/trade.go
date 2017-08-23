package controller

import (
	"log"

	"github.com/gin-gonic/gin"
)

//trade manage

var (
	trade *gin.RouterGroup
)

func init() {
	if Server == nil {
		log.Fatalln("trade:init web server error")
		return
	}
	//init trade group tp:type op:oprate
	trade = Server.Group("/trade", gin.Logger(), gin.Recovery(), getToken(), AuthCheck())
}
