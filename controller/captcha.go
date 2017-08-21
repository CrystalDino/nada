package controller

import (
	"bytes"
	"log"
	"nada/core"
	"net/http"
	"path"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func init() {
	if Server == nil {
		log.Fatalln("captcha:init web server error")
		return
	}
	Server.GET("/captcha", DealCaptcha)
}

func DealCaptcha(c *gin.Context) {
	tp, has := c.GetQuery("type")
	if !has {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	switch tp {
	case "id":
		makeCaptchaId(c)
	case "pic":
		getCaptcha(c)
	// case "chk":
	// 	checkCaptcha(c)
	default:
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func checkCaptcha(c *gin.Context) {
	k, has := c.GetQuery("key")
	if !has || k == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	v, has := c.GetQuery("value")
	if !has || v == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	r := NewResult()
	if captcha.VerifyString(k, v) {
		r["Ok"] = true
	}
	c.JSON(http.StatusOK, r)
}

func makeCaptchaId(c *gin.Context) {
	r := NewResult()
	id := captcha.New()
	r["Id"] = id
	r["Ok"] = true
	c.JSON(http.StatusOK, r)
}

func getCaptcha(c *gin.Context) {
	idvalue, has := c.GetQuery("name")
	if !has || idvalue == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, file := path.Split(idvalue)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext != ".png" || id == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if c.Query("reload") != "" {
		captcha.Reload(id)
	}

	var content bytes.Buffer
	width, height := core.GlobalConfig.GetCaptchaSize()
	if err := captcha.WriteImage(&content, id, width, height); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Writer.Header().Set("Pragma", "no-cache")
	c.Writer.Header().Set("Expires", "0")
	c.Data(http.StatusOK, "image/png", content.Bytes())
}
