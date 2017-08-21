package controller

//asset manage

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	asset *gin.RouterGroup
)

func init() {
	if Server == nil {
		log.Fatalln("asset:init web server error")
		return
	}
	//init asset group tp:type op:oprate
	asset = Server.Group("/asset", getToken(), AuthCheck())
	asset.GET("/recharge", RechargeLog)
	asset.POST("/recharge/:op", Recharge)

	asset.GET("/withdraw/:tp", WithdrawLog)
	asset.POST("/withdraw/:tp", Withdraw)

	asset.GET("/coinin/:id/:tp", CoinInLog)
	asset.POST("/coinin/:id/:op", CoinIn)

	asset.GET("/coinout/:id/:tp", CoinOutLog)
	asset.POST("/coinout/:id/:op", CoinOut)

	asset.GET("/bankcard", BankCardInfo)
	asset.POST("/bankcard/:op", BankCard)

	asset.GET("/trust/:id/:tp", TrustLog)
	asset.POST("/trust/:id/:op", Trust)

	asset.GET("/order/:id/:tp", OrderLog)
}

//RechargeLog get recharge log
func RechargeLog(c *gin.Context) {
	c.String(http.StatusOK, "recharge log")
}

//Recharge make a new recharge order
func Recharge(c *gin.Context) {
	r := NewResult()
	r["Err"] = "not support"
	c.JSON(http.StatusOK, r)
}

//WithdrawLog get withdraw log
func WithdrawLog(c *gin.Context) {
	c.String(http.StatusOK, "withdraw log")
}

//Withdraw name a new withdraw order
func Withdraw(c *gin.Context) {
	r := NewResult()
	r["Err"] = "not support"
	c.JSON(http.StatusOK, r)
}

//CoinIn recharge virtual coin, apply a virtual coin address
func CoinIn(c *gin.Context) {
	c.String(http.StatusOK, "coinin %s-%s", c.Param("id"), c.Param("op"))
}

//CoinInLog show virtual coin rechage log
func CoinInLog(c *gin.Context) {
	c.String(http.StatusOK, "coinin log %s", c.Param("id"))
}

//CoinOut withdraw virtual coin, accept a virtual coin address
func CoinOut(c *gin.Context) {
	c.String(http.StatusOK, "coinout %s", c.Param("id"))
}

//CoinOutLog show virtual coin withdraw log
func CoinOutLog(c *gin.Context) {
	c.String(http.StatusOK, "coinout log %s", c.Param("id"))
}

//DealLog show deal log
func OrderLog(c *gin.Context) {
	c.String(http.StatusOK, "deal log %s-%s", c.Param("id"), c.Param("tp"))
}

//Trust make a new order buy&sell
func Trust(c *gin.Context) {}

//TrustLog show order log that have not deal
func TrustLog(c *gin.Context) {}

//BankCard add or remove a bank card
func BankCard(c *gin.Context) {}

//BankCardInfo show bank cards info
func BankCardInfo(c *gin.Context) {}
