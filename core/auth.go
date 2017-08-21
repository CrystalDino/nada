package core

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func transToClaims(data interface{}, except ...string) (mc jwt.MapClaims, err error) {
	m, err := StructToMap(data, true, except...)
	if err != nil {
		return
	}
	mc = jwt.MapClaims(m)
	return
}

func transToData(mc jwt.MapClaims) (data map[string]interface{}) {
	return map[string]interface{}(mc)
}

//Unmarshal todo
func Unmarshal(src map[string]interface{}, dst interface{}) (err error) {
	// if dst == nil || len(src) == 0 {
	// 	return errors.New("dst or src is nil")
	// }
	// dv := reflect.ValueOf(dst)
	// if dv.Kind() != reflect.Ptr {
	// 	return errors.New("type of dst must be ptr")
	// }

	return
}

func TokenMake(data interface{}, except ...string) (tokenStr string, err error) {
	var mc jwt.MapClaims
	mc, err = transToClaims(data, except...)
	if err != nil {
		return
	}
	mc["expire"] = time.Now().Unix() + GlobalConfig.GetExpire()
	tokenStr, err = jwt.NewWithClaims(jwt.SigningMethodHS256, mc).SignedString(GlobalConfig.GetTokenKey())
	return
}

func TokenValidate(tokenStr string) (data interface{}, err error) {
	if tokenStr == "" {
		err = errors.New("token is nil")
		return
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return GlobalConfig.GetTokenKey(), nil
	})
	if err != nil {
		return
	}
	if !token.Valid {
		err = jwt.ErrInvalidKey
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("get claims from token error")
		return
	}
	expire, ok := claims["expire"]
	if !ok {
		err = errors.New("token lost expire")
		return
	}
	if int64(expire.(float64)) < time.Now().Unix() {
		err = errors.New("token expire")
		return
	}
	data = transToData(claims)
	return

}
