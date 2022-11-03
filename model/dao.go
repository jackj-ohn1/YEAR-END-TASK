package model

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
)

var D *Database = &Database{DB: new(gorm.DB)}

type Database struct {
	DB *gorm.DB
}

func (u *UserModel) AddUser() error {
	return D.DB.Table("users").Save(u).Error
}

func (c *ConsumeModel) AddRecord(tx *gorm.DB) error {
	return tx.Table("card_records").Create(c).Error
}

func (l *LibraryModel) AddBook(tx *gorm.DB) error {
	return tx.Table("books").Create(l).Error
}

func GetAllUsers() ([]string, error) {
	var users []string
	err := D.DB.Table("users").Pluck("id", &users).Error
	return users, err
}

func GetTotalMoney(uname string) float64 {
	var sum []float64
	D.DB.Table("card_records").
		Where("id=?", uname).
		Pluck("SUM(trans_money)", &sum)
	ret, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", sum[0]), 64)
	return ret
}

func GetTotalRecord(uname string) int {
	var count int
	D.DB.Table("card_records").
		Where("id=?", uname).
		Count(&count)
	return count
}

func GetLocationTime(uname string) []OrderData {
	var location []OrderData
	
	D.DB.Table("card_records").
		Group("location").
		Select([]string{"location as name", "COUNT(location) as count"}).
		Where("id=?", uname).
		Order("count desc").
		Scan(&location)
	
	return location
}

func GetMethodTime(uname string) []OrderData {
	var method []OrderData
	
	D.DB.Table("card_records").
		Group("method").
		Select([]string{"method as name", "COUNT(method) as count"}).
		Where("id=?", uname).
		Order("count desc").
		Scan(&method)
	
	return method
}

func GetBooks(uname string) []string {
	var names []string
	D.DB.Table("library").
		Where("id=?", uname).
		Pluck("name", &names)
	return names
}

func (d *Database) Init() {
	var err error
	d.DB, err = gorm.Open("mysql", "yyj:0118@/year_end?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
}

func Start(ctx context.Context) {
	D.Init()
	DR.Init()
	SetAllUserRecord(ctx)
}
