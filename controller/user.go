package controller

//user manage

import (
	"log"
	"nada/core"
	"nada/models"
	"net/http"
	"time"

	"reflect"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

var (
	user *gin.RouterGroup
)

func init() {
	if Server == nil {
		log.Fatalln("user:init web server error")
		return
	}
	Server.POST("/user/login", VerifyCodeCheck(), UserLogin)
	Server.POST("/user/register", VerifyCodeCheck(), UserRegister)
	//init user group
	user = Server.Group("/user", gin.Logger(), gin.Recovery(), getToken(), AuthCheck())
	user.GET("/logout", UserLogout)
	user.GET("/info", UserInfo)
	user.POST("/updatepwd", UpdatePasswd)
	user.POST("/updatetsc", UpdateTranscode)
	user.POST("/info", UpdateInfo)
	user.POST("/identification", UserIdentificate)
}

//UserLogin login method
func UserLogin(c *gin.Context) {
	r := core.NewResult()
	ufl := &models.UserForLogin{}
	if err := c.Bind(ufl); err != nil {
		println(err.Error())
		return
	}
	//check password
	u, err := ufl.UserPasswdCheck()
	if err != nil {
		log.Println(err)
		r.SetErr("login failed")
		c.JSON(http.StatusOK, r)
		return
	}

	//make token
	if token, err := u.CreateToken(); err != nil {
		r.SetErr(err.Error())
	} else {
		r.Set("Nada", token)
		r.SetOk(true)
	}

	//update user table loginip & login time
	u.UpdateLoginInfo(c.ClientIP(), time.Now().Unix())

	c.JSON(http.StatusOK, r)
}

//UserRegister user register
func UserRegister(c *gin.Context) {
	ufr := &models.UserForRegister{}
	r := core.NewResult()
	err := c.Bind(ufr)
	if err != nil {
		return
	}
	u, err := ufr.ToUser()
	if err != nil {
		r.SetErr("internal error")
		c.JSON(http.StatusOK, r)
		return
	}
	u.LastLoginIp = c.ClientIP()
	u.LTime = time.Now().Unix()
	id, err := u.Stor()
	if err != nil {
		r.SetErr(err.Error())
	} else {
		r.Set("Id", id)
		r.SetOk(true)
	}
	c.JSON(http.StatusOK, r)
}

//UserLogout logout method，delete token at client/browse
func UserLogout(c *gin.Context) {
	c.String(http.StatusOK, "%s", "user logout")
}

//UserInfo show detail info of user
func UserInfo(c *gin.Context) {
	r := core.NewResult()
	u, err := models.GetUserByID(c.GetInt64("uid"))
	if err != nil {
		r.SetErr(err.Error())
	} else {
		r.Set("User", u)
		r.SetOk(true)
	}
	c.JSON(http.StatusOK, r)
}

//UpdatePasswd update user password
func UpdatePasswd(c *gin.Context) {
	r := core.NewResult()
	oPwd, has := c.GetPostForm("oPwd")
	if !has {
		r.SetErr("lost oPwd")
		c.JSON(http.StatusOK, r)
		return
	}
	nPwd, has := c.GetPostForm("nPwd")
	if !has {
		r.SetErr("lost nPwd")
		c.JSON(http.StatusOK, r)
		return
	}
	rPwd, has := c.GetPostForm("rPwd")
	if !has {
		r.SetErr("lost rPwd")
		c.JSON(http.StatusOK, r)
		return
	}
	if oPwd == "" || nPwd == "" || rPwd == "" {
		r.SetErr("pwd can not be nil")
		c.JSON(http.StatusOK, r)
		return
	}
	if nPwd != rPwd {
		r.SetErr("nPwd not equal to rPwd")
		c.JSON(http.StatusOK, r)
		return
	}
	err := models.UpdatePassword(c.GetInt64("uid"), oPwd, nPwd)
	if err != nil {
		r.SetErr(err.Error())
		c.JSON(http.StatusOK, r)
		return
	}
	r.SetOk(true)
	c.JSON(http.StatusOK, r)
}

//UpdateTranscode update transacton password
func UpdateTranscode(c *gin.Context) {
	r := core.NewResult()
	oPwd, has := c.GetPostForm("oPwd")
	if !has {
		r.SetErr("lost oPwd")
		c.JSON(http.StatusOK, r)
		return
	}
	nPwd, has := c.GetPostForm("nPwd")
	if !has {
		r.SetErr("lost nPwd")
		c.JSON(http.StatusOK, r)
		return
	}
	rPwd, has := c.GetPostForm("rPwd")
	if !has {
		r.SetErr("lost rPwd")
		c.JSON(http.StatusOK, r)
		return
	}
	if nPwd == "" || rPwd == "" {
		r.SetErr("pwd can not be nil")
		c.JSON(http.StatusOK, r)
		return
	}
	if nPwd != rPwd {
		r.SetErr("nPwd not equal to rPwd")
		c.JSON(http.StatusOK, r)
		return
	}
	err := models.UpdateTranscode(c.GetInt64("uid"), oPwd, nPwd)
	if err != nil {
		r.SetErr(err.Error())
		c.JSON(http.StatusOK, r)
		return
	}
	r.SetOk(true)
	c.JSON(http.StatusOK, r)
}

//UpdateInfo update user info except password and transaction password
func UpdateInfo(c *gin.Context) {
	ufi := &models.UserForUpdateInfo{}
	if err := c.Bind(ufi); err != nil {
		return
	}
	r := core.NewResult()
	err := ufi.Update(c.GetInt64("uid"))
	if err != nil {
		r.SetErr(err.Error())
		c.JSON(http.StatusOK, r)
		return
	}
	r.SetOk(true)
	c.JSON(http.StatusOK, r)
}

//UserIdentificate identification of user
func UserIdentificate(c *gin.Context) {
	c.String(http.StatusOK, "user identificate")
}

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		tv := c.MustGet(core.DefaultInternalTokenName).(string)
		data, err := core.TokenValidate(tv)
		if err != nil {
			r := core.MakeResult(false, err.Error())
			c.AbortWithStatusJSON(http.StatusNotAcceptable, r)
			return
		}
		if v, ok := data.(map[string]interface{})["id"]; !ok {
			r := core.MakeResult(false, "invalid token no id")
			c.AbortWithStatusJSON(http.StatusNotAcceptable, r)
			return
		} else {
			if reflect.TypeOf(v).Kind() != reflect.Float64 {
				r := core.MakeResult(false, "invalid token wrong id type")
				c.AbortWithStatusJSON(http.StatusNotAcceptable, r)
				return
			}
			c.Set("uid", int64(v.(float64)))
			c.Set(core.GlobalConfig.GetTokenName(), data)
		}
		c.Next()
	}
}

func VerifyCodeCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, has := c.GetPostForm("id")
		if !has {
			r := core.MakeResult(false, "verify params lost id")
			c.AbortWithStatusJSON(http.StatusBadRequest, r)
			return
		}
		code, has := c.GetPostForm("code")
		if !has {
			r := core.MakeResult(false, "verify params lost code")
			c.AbortWithStatusJSON(http.StatusBadRequest, r)
			return
		}
		if id == "" || code == "" {
			r := core.MakeResult(false, "verify params is empty")
			c.AbortWithStatusJSON(http.StatusOK, r)
			return
		}
		if !captcha.VerifyString(id, code) {
			r := core.MakeResult(false, "verify failed")
			c.AbortWithStatusJSON(http.StatusOK, r)
			return
		}
		c.Next()
	}
}
