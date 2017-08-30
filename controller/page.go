package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	if Server == nil {
		log.Fatalln("page:init web server error")
		return
	}

	Server.GET("/", Home)
	Server.GET("/login", Login)
	Server.GET("/register", Register)

}

//Home home page
func Home(c *gin.Context) {
	c.HTML(http.StatusOK,"home.html",gin.H{
		"title":"crystal.dino",
	})
}

//Login the login page
func Login(c *gin.Context) {
	c.String(http.StatusOK, "login")
}

//Register user register page
func Register(c *gin.Context) {
	c.String(http.StatusOK, "resiger")
}
