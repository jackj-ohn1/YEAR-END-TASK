package model

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

// 数据库model

type UserModel struct {
	ID       string `gorm:"id"`
	Password string `gorm:"password"`
}

type ConsumeModel struct {
	ID         string  `gorm:"id"`
	TransMoney float64 `gorm:"tans_money"`
	Location   string  `gorm:"location"`
	Date       string  `gorm:"data"`
	Method     string  `gorm:"method"`
}

type LibraryModel struct {
	Name string `gorm:"name"`
	Id   string `gorm:"id"`
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
	Total     int32           `json:"total"`
	ViewEnd   int32           `json:"viewEnd"`
	Records   int32           `json:"records"`
	Rows      []ConsumeRecord `json:"rows"`
	Start     int32           `json:"start"`
	PageSize  int32           `json:"pageSize"`
	ViewStart int32           `json:"viewStart"`
	Page      int32           `json:"page"`
}

// 支付方式 支付地点 支付金额
type ConsumeRecord struct {
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

func (record *ConsumeRecord) ConvertToModel(uname string) ConsumeModel {
	return ConsumeModel{
		ID:         uname,
		TransMoney: record.TransMoney,
		Location:   strings.Replace(record.OrgName, "华中师范大学/后勤集团/", "", -1),
		Method:     record.WalletName,
		Date:       record.DealDate,
	}
}
