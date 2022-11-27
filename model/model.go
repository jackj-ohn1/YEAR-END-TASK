package model

import (
	jsoniter "github.com/json-iterator/go"
	_ "gorm.io/driver/mysql"
	"strings"
)

// 数据库model

type User struct {
	Id       uint   `gorm:"column:id;type:int;autoIncrement;primaryKey"`
	Account  string `gorm:"column:account;type:varchar(20);primaryKey"`
	Password string `gorm:"column:password;type:varchar(30);not null;"`
}

type CardRecord struct {
	ID         uint    `gorm:"column:id;type:int;autoIncrement;primaryKey"`
	Account    string  `gorm:"column:account;type:varchar(20)"`
	TransMoney float64 `gorm:"column:trans_money;type:float"`
	Location   string  `gorm:"column:location;type:varchar(100)"`
	Date       string  `gorm:"column:data;type:varchar(100)"`
	Method     string  `gorm:"column:method;type:varchar(100)"`
}

type Book struct {
	ID      uint   `gorm:"id;type:int;primaryKey;autoIncrement"`
	Name    string `gorm:"name;type:varchar(30)"`
	Account string `gorm:"account;type:varchar(20)"`
}

//

type OrderData struct {
	Name  string `gorm:"name"`
	Count int    `gorm:"count"`
}

func (o OrderData) MarshalBinary() (data []byte, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(&o)
}

func (o OrderData) UnmarshalBinary(data []byte) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, &o)
}

// 数据模板
type Record struct {
	Errcode string `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Result  rows   `json:"result"`
}

type rows struct {
	Total     int32               `json:"total"`
	ViewEnd   int32               `json:"viewEnd"`
	Records   int32               `json:"records"`
	Rows      []SelfConsumeRecord `json:"rows"`
	Start     int32               `json:"start"`
	PageSize  int32               `json:"pageSize"`
	ViewStart int32               `json:"viewStart"`
	Page      int32               `json:"page"`
}

// 支付方式 支付地点 支付金额
type SelfConsumeRecord struct {
	DealName   string  `json:"dealName"`
	OrgName    string  `json:"orgName"`
	DealDate   string  `json:"dealDate"'`
	TransNo    string  `json:"-"`
	TransMoney float64 `json:"transMoney"`
	WalletName string  `json:"walletName"`
	CardNo     string  `json:"cardNo"`
	InMoney    float64 `json:"inMoney"`
	OutMoney   float64 `json:"outMoney"`
}

func (record *SelfConsumeRecord) ConvertToModel(uname string) CardRecord {
	return CardRecord{
		Account:    uname,
		TransMoney: record.TransMoney,
		Location:   strings.Replace(record.OrgName, "华中师范大学/后勤集团/", "", -1),
		Method:     record.WalletName,
		Date:       record.DealDate,
	}
}
