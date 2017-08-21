package models

type Order struct {
	Id     int64
	Buid   int64   `xorm:"notnull defualt(0)"`                 //buyer user id
	Suid   int64   `xorm:"notnull default(0)"`                 //sale user id
	Btid   int64   `xorm:"notnull default(0)"`                 //trust buy id
	Stid   int64   `xorm:"notnull default(0)"`                 //trust sale id
	Brate  float64 `xorm:"decimal(10,5) notnull default(0.0)"` //buy rate
	Srate  float64 `xorm:"decimal(10,5) notnull default(0.0)"` //sale rate
	Price  float64 `xorm:"decimal(20,8) notnull default(0.0)"`
	Amount float64 `xorm:"decimal(20,8) notnull default(0.0)"`
	Ctime  int64   `xorm:"created"`
	Tp     int8    `xorm:"tinyint(1) notnull default(-1)"` //成交类型 1:sale 0:buy
}

func (order *Order) TableName() string {
	return "bargain"
}
