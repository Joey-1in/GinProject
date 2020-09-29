package model

type Menber struct {
	Id int64 `xorm:"pk autoincr" json:"id"`
	UserName string `xorm:"varchar(20)" json:"userName"`
	Mobile string `xorm:"varchar(11)" json:"mobile"`
	Password string `xorm:"varchar(255)" json:"password"`
	RegisterTime int64 `xorm:"bigint" json:"registerTime"`
	Avatar string `xorm:"varchar(255)" json:"avatar"`
	Balance float64 `xorm:"double" json:"balance"`
	Isactive int8 `xorm:"tinyint" json:"isactive"`
	City string `xorm:"varchar(10)" json:"city"`
}