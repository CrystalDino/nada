package models

import (
	"log"
	"nada/core"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	DSN    string
	engine *xorm.Engine
)

func init() {
	var err error
	DSN = core.GlobalConfig.GetDSN()
	if engine, err = xorm.NewEngine("mysql", DSN); err != nil {
		log.Fatalln("link to DB error,", err)
		return
	}

	if err = SyncModels(); err != nil {
		log.Fatalln("sync db failed,", err)
		return
	}

	log.Println("make new engine done--------------------------------------")
}

func SyncModels() (err error) {
	if err = engine.Sync2(&User{}); err != nil {
		return
	}

	if err = engine.Sync2(&Order{}); err != nil {
		return
	}
	return nil
}
