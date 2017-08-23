package controller

//user manage

import (
	"log"
	"nada/core"
	"nada/models"
	"net/http"
	"time"

	"reflect"

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
	Server.POST("/user/login", UserLogin)
	Server.POST("/user/register", UserRegister)
	//init user group
	user = Server.Group("/user", getToken(), AuthCheck())
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
		return
	}

	//check verify code
	// if ufl.CheckCode == "" || ufl.CheckID == "" {
	// 	r["Err"] = "field lost"
	// 	c.JSON(http.StatusOK, r)
	// 	return
	// }
	// if !captcha.VerifyString(ufl.CheckID, ufl.CheckCode) {
	// 	r["Err"] = "verify failed"
	// 	c.JSON(http.StatusOK, r)
	// 	return
	// }

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
	u := &models.User{}
	err := c.Bind(u)
	if err != nil {
		return
	}
	u.LastLoginIp = c.ClientIP()
	u.LTime = time.Now().Unix()
	r := core.NewResult()
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
