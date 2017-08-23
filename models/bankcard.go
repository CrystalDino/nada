package models

import (
	"errors"
	"time"
)

type BankCard struct {
	Id       int64
	UserId   int64  `xorm:"notnull default 0"`
	Status   int8   `xorm:"tinyint(1) notnull default 0"` //1:冻结 0：开通
	CardType int8   `xorm:"tinyint(1) notnull default 0"` //卡片类型 1借记卡 2信用卡
	CTime    int64  `xorm:"created notnull"`
	MTime    int64  `xorm:"updated notnull"`
	CardNo   string `xorm:"char(20) index unique notnull default ''"`
	BankAddr string `xorm:"varchar(256) notnull default ''"`
	Remark   string `xorm:"varchar(128) notnull default ''"`
	BankName string `xorm:"varchar(64) notnull default ''"`
}

type BankCardForAdd struct {
	CardNo   string `form:"cardno" binding:"required"`
	CardType int8   `form:"cardtype" binding:"required"`
	BankAddr string `form:"bankaddr" binding:"required"`
	Remark   string `form:"remark" binding:"required"`
	BankName string `form:"bankname" binding:"required"`
}

func (bc *BankCard) TableName() string {
	return "bankcard"
}

func (bca *BankCardForAdd) ToBankCard(uid int64) (bc *BankCard) {
	bc = &BankCard{
		CardNo:   bca.CardNo,
		CardType: bca.CardType,
		BankAddr: bca.BankAddr,
		Remark:   bca.Remark,
		BankName: bca.BankName,
		UserId:   uid,
		CTime:    time.Now().Unix(),
		MTime:    time.Now().Unix(),
	}
	//todo: add some check here

	return
}

func (bc *BankCard) Stor() (id int64, err error) {
	if bc.CardType < 0 || bc.CardType > 2 {
		return -1, errors.New("card type error")
	}
	if id, err = engine.InsertOne(bc); err != nil {
		return -1, err
	}
	if id != 1 {
		return -1, errors.New("db error: insert count not one")
	}
	var cbc = &BankCard{CardNo: bc.CardNo}
	has, err := engine.Select("id").Get(cbc)
	if err != nil {
		return -1, err
	}
	if !has {
		return -1, errors.New("create bank card record " + bc.CardNo + " failed")
	}
	// log.Println("password:", u.Password)
	id = cbc.Id
	return
}
