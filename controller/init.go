package controller

import (
	"log"
	"nada/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	//Server web server engine - *gin.Engine
	Server = gin.New()
)

func init() {
	// Server = gin.New()
	if Server == nil {
		log.Fatalln("init web server error")
		return
	}
	// Global middleware
	Server.Use(gin.Logger())
	Server.Use(gin.Recovery())
	//add view
	Server.StaticFS("/static", http.Dir("view/static"))
	//load template
	Server.LoadHTMLGlob("view/templates/*")
}

func getToken() gin.HandlerFunc {
	tn := core.GlobalConfig.GetTokenName()
	return func(c *gin.Context) {
		tv := c.GetHeader(tn)
		if tv != "" {
			c.Set(core.DefaultInternalTokenName, tv)
			c.Next()
			return
		}
		var has bool
		if tv, has = c.GetQuery(tn); has {
			c.Set(core.DefaultInternalTokenName, tv)
			c.Next()
			return
		}
		r := core.MakeResult(false, "no token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, r)
	}
}

func Run() {
	if Server == nil {
		log.Fatalln("web server not running")
		return
	}
	Server.StaticFile("/test/jwt", "./test/jwt.html")
	Server.Run(core.GlobalConfig.GetServerAddr())
}
