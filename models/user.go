package models

import (
	"errors"
	"log"
	"nada/core"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id          int64
	Name        string `xorm:"char(32) notnull" form:"username" binding:"required"`
	Cell        string `xorm:"char(16) notnull index unique" form:"cell" binding:"required"`
	Password    string `xorm:"char(64) notnull" form:"password" binding:"required"`
	Transcode   string `xorm:"char(64)"`
	Email       string `xorm:"char(32)" form:"email" binding:"required"`
	LastLoginIp string `xorm:"char(16)"`
	Stat        int8
	LTime       int64
	CTime       int64  `xorm:"created notnull"`
	MTime       int64  `xorm:"updated notnull"`
	DTime       int64  `xorm:"deleted"`
	RealName    string `xorm:"char(8)"`
	Sex         int8   `xorm:"tinyint(1)"`
	CardNo      string `xorm:"char(20)"`
	Area        string `xorm:"varchar(128)"`
	Icon        string `xorm:"char(128)" form:"icon"`
	Info        string `xorm:"varchar(256)" form:"info"`
}

type UserForLogin struct {
	Cell      string `form:"cell" binding:"required"`
	Password  string `form:"password" binding:"required"`
	CheckCode string `form:"code" binding:"required"`
	CheckID   string `form:"id" binding:"required"`
}

type UserForUpdateInfo struct {
	Name  string `form:"username"`
	Email string `form:"email"`
	Icon  string `form:"icon"`
	Info  string `form:"info"`
}

type UserForIdentificate struct {
	RealName    string `form:"name"`
	Sex         string `form:"sex"`
	CardNo      string `form:"no"`
	Area        string `form:"area"`
	CardFront   string `form:"front"`
	CardRear    string `form:"rear"`
	CardInHands string `form:"inhands"`
}

func (user *User) TableName() string {
	return "user"
}

func (ufl *UserForLogin) UserPasswdCheck() (u *User, err error) {
	if ufl.Cell == "" || ufl.Password == "" {
		err = errors.New("params lost")
		return
	}
	u = &User{Cell: ufl.Cell}
	has, err := engine.Get(u)
	if err != nil {
		u = nil
		return
	}
	if !has {
		u, err = nil, errors.New("no this user")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(ufl.Password))
	if err != nil {
		u = nil
		return
	}
	return
}

func (u *User) CreateToken() (token string, err error) {
	token, err = core.TokenMake(u, "dtime", "password", "transcode", "stat", "ctime",
		"mtime", "icon", "info", "lastloginip", "ltime", "email", "cardno", "area", "realname", "sex")
	return
}

func (u *User) Stor() (id int64, err error) {
	pas, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}
	u.Stat = 1 //new user
	u.Password = string(pas)
	if id, err = engine.InsertOne(u); err != nil {
		return -1, err
	}
	if id != 1 {
		return -1, errors.New("db error: insert count not one")
	}
	var cu = &User{Cell: u.Cell}
	has, err := engine.Select("id").Get(cu)
	if err != nil {
		return -1, err
	}
	if !has {
		return -1, errors.New("create user account " + u.Name + " failed")
	}
	// log.Println("password:", u.Password)
	id = cu.Id
	return
}

func GetUserByID(id int64) (m map[string]interface{}, err error) {
	u, m := &User{}, make(map[string]interface{})
	if id <= 0 {
		err = errors.New("user id error")
		return
	}
	has, err := engine.Id(id).Cols("id", "name", "cell", "email", "last_login_ip", "l_time", "c_time", "m_time", "icon", "info").Get(u)
	if err != nil {
		log.Println("get user by id fail,", err)
		err = errors.New("get info fail")
		return
	}
	if !has {
		log.Println("can not find user by id,", id)
		err = errors.New("no user data")
		return
	}
	m, err = core.StructToMap(u, false, "Password", "Transcode", "DTime", "Stat")
	if err != nil {
		log.Println("get user by id fail,", err)
		err = errors.New("user to map fail")
		return
	}
	return
}

func UpdatePassword(id int64, oPwd, nPwd string) (err error) {
	if id <= 0 {
		return errors.New("user id error")
	}
	if !core.GlobalConfig.CheckPwdLength(oPwd) {
		return errors.New("oPwd length error")
	}
	if !core.GlobalConfig.CheckPwdLength(nPwd) {
		return errors.New("nPwd length error")
	}
	u := &User{}
	has, err := engine.Id(id).Cols("id", "password").Get(u)
	if err != nil {
		log.Println("get user by id fail for updating pwd,", err)
		err = errors.New("get info fail")
		return
	}
	if !has {
		log.Println("can not find user by id for updating pwd,", id)
		err = errors.New("no user data")
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oPwd)); err != nil {
		return
	}
	pas, err := bcrypt.GenerateFromPassword([]byte(nPwd), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	u.Password = string(pas)
	n, err := engine.Id(u.Id).Cols("password").Update(u)
	if err != nil {
		return
	}
	if n != 1 {
		log.Println("update user password for id=", id, "failed, by db row affect=", n)
		return errors.New("update not one")
	}
	return
}

func UpdateTranscode(id int64, oPwd, nPwd string) (err error) {
	if id <= 0 {
		return errors.New("user id error")
	}
	if !core.GlobalConfig.CheckPwdLength(oPwd) {
		return errors.New("oPwd length error")
	}
	if !core.GlobalConfig.CheckPwdLength(nPwd) {
		return errors.New("nPwd length error")
	}
	u := &User{}
	has, err := engine.Id(id).Cols("id", "transcode").Get(u)
	if err != nil {
		log.Println("get user by id fail for updating tsc,", err)
		err = errors.New("get info fail")
		return
	}
	if !has {
		log.Println("can not find user by id for updating tsc,", id)
		err = errors.New("no user data")
		return
	}
	if len(u.Transcode) != 0 {
		if oPwd == "" {
			return errors.New("oPwd is needed")
		}
		if err = bcrypt.CompareHashAndPassword([]byte(u.Transcode), []byte(oPwd)); err != nil {
			return
		}
	}

	pas, err := bcrypt.GenerateFromPassword([]byte(nPwd), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	u.Transcode = string(pas)
	n, err := engine.Id(u.Id).Cols("transcode").Update(u)
	if err != nil {
		return
	}
	if n != 1 {
		log.Println("update user transcode for id=", id, "failed, by db row affect=", n)
		return errors.New("update not one")
	}
	return
}

func (u *User) UpdateLoginInfo(ip string, tm int64) (err error) {
	u.LastLoginIp = ip
	u.LTime = tm
	n, err := engine.Id(u.Id).Cols("last_login_ip", "ltime").Update(u)
	if err != nil {
		return
	}
	if n != 1 {
		err = errors.New("update not one")
	}
	return
}

func (ufi *UserForUpdateInfo) Update(id int64) (err error) {
	if id <= 0 {
		return errors.New("user id error")
	}
	u := &User{
		Id:    id,
		Name:  ufi.Name,
		Email: ufi.Email,
		Icon:  ufi.Icon,
		Info:  ufi.Info,
	}
	n, err := engine.Id(id).Cols("name", "email", "icon", "info").Update(u)
	if err != nil {
		return
	}
	if n != 1 {
		err = errors.New("update not one")
	}
	return
}
