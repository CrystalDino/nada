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
	Server.GET("/news", News)
	Server.GET("/summary", Summary)
	Server.GET("/exchange", Exchange)
	Server.GET("/account", Account)
	Server.GET("/assets", Assets)
	Server.GET("/code", Code)
}

//Home home page
func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title":  "Home",
		"isHome": true,
	})
}

//News news page
func News(c *gin.Context) {
	c.HTML(http.StatusOK, "news.html", gin.H{
		"title":  "News",
		"isNews": true,
	})
}

//Summary summary page
func Summary(c *gin.Context) {
	c.HTML(http.StatusOK, "summary.html", gin.H{
		"title":     "Summary",
		"isSummary": true,
	})
}

//Exchage exchange page
func Exchange(c *gin.Context) {
	c.HTML(http.StatusOK, "exchange.html", gin.H{
		"title":      "Exchange",
		"isExchange": true,
	})
}

//Account account page
func Account(c *gin.Context) {
	c.HTML(http.StatusOK, "account.html", gin.H{
		"title":     "Account",
		"isAccount": true,
	})
}

//Assets assets page
func Assets(c *gin.Context) {
	c.HTML(http.StatusOK, "assets.html", gin.H{
		"title":    "Assets",
		"isAssets": true,
	})
}

//Code github page
func Code(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "https://github.com/CrystalDino")
}

//Login the login page
func Login(c *gin.Context) {
	c.String(http.StatusOK, "login")
}

//Register user register page
func Register(c *gin.Context) {
	c.String(http.StatusOK, "resiger")
}
