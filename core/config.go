package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

const (
	DefaultTokenKey          string = "crystal.dino"
	DefaultServerAddr        string = ":9630"
	DefaultTokenName         string = "Nada"
	DefaultInternalTokenName string = "token"
	DefaultTokenExpire       int64  = 3600
	DefaultDSN               string = "dino:passw0rd@tcp(localhost:3306)/valet?charset=utf8"
	DefaultCaptchaWidth      int    = 240
	DefaultCaptchaHeight     int    = 80
	DefaultPwdMin            int    = 8
	DefaultPwdMax            int    = 16
)

var (
	GlobalConfig Config
)

type Config struct {
	TokenKey    string `json:"token_key"`
	TokenExpire int64  `json:"token_expire"`

	LogInfoPath  string `json:"log_info"`
	LogErrorPath string `json:"log_error"`

	DSN string `json:"db_source"`

	ServerAddr string `json:"server_addr"`
	TokenName  string `json:"token_name"`

	CaptchaWidth  int `json:"captcha_width"`
	CaptchaHeight int `json:"captcha_height"`

	PwdMin int `json:"pwd_min"`
	PwdMax int `json:"pwd_max"`
}

func LoadConfigFile(filePath string) (err error) {
	//当热更新配置的时候，应该需要个锁
	if filePath == "" {
		err = errors.New("the path of config file is nil")
		return
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	if err = json.Unmarshal(data, &GlobalConfig); err != nil {
		return
	}
	//参数检测
	if err = GlobalConfig.ValidCheck(); err != nil {
		return
	}

	return
}

func (c *Config) ValidCheck() (err error) {
	if c.PwdMin < 0 || c.PwdMax < 0 || c.PwdMin > c.PwdMax {
		return errors.New("length of pwd setup error")
	}
	if c.CaptchaHeight < 0 || c.CaptchaWidth < 0 {
		return errors.New("size of captcha setup error")
	}
	return
}

func (c *Config) CheckPwdLength(pwd string) bool {
	min, max := DefaultPwdMin, DefaultPwdMax
	if c.PwdMin != 0 {
		min = c.PwdMin
	}
	if c.PwdMax != 0 {
		max = c.PwdMax
	}
	if len(pwd) < min || len(pwd) > max {
		return false
	}
	return true
}

func (c *Config) GetCaptchaSize() (width, height int) {
	if c.CaptchaWidth == 0 {
		width = DefaultCaptchaWidth
	} else {
		width = c.CaptchaWidth
	}
	if c.CaptchaHeight == 0 {
		height = DefaultCaptchaHeight
	} else {
		height = c.CaptchaHeight
	}
	return
}

func (c *Config) GetDSN() string {
	if c.DSN == "" {
		return DefaultDSN
	}
	return c.DSN
}

func (c *Config) GetTokenKey() []byte {
	if c.TokenKey == "" {
		return []byte(DefaultTokenKey)
	}
	return []byte(c.TokenKey)
}

func (c *Config) GetServerAddr() string {
	if c.ServerAddr == "" {
		return DefaultServerAddr
	}
	return c.ServerAddr
}

func (c *Config) GetTokenName() string {
	if c.TokenName == "" {
		return DefaultTokenName
	}
	return c.TokenName
}

func (c *Config) GetExpire() int64 {
	if c.TokenExpire == 0 {
		return DefaultTokenExpire
	}
	return c.TokenExpire
}
