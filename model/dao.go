package model

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"
	"year-end/utils/errno"
)

var D *Database = &Database{DB: new(gorm.DB)}

type Database struct {
	DB *gorm.DB
}

func (u *User) CheckUser() error {
	var user User
	if err := D.DB.Table("users").Where("account=?", u.Account).Find(&user).Error; err != nil {
		return err
	}
	if user.Account == u.Account {
		return errors.New(errno.GetMessage(errno.ERR_USER_EXIST))
	}
	// 可以插入
	return nil
}

func (u *User) AddUser() error {
	if err := u.CheckUser(); err != nil {
		return err
	}
	return D.DB.Table("users").Create(u).Error
}

func (c *CardRecord) AddRecord(tx *gorm.DB) error {
	return tx.Table("card_records").Omit("id").Create(c).Error
}

func (l *Book) AddBook(tx *gorm.DB) error {
	return tx.Table("books").Create(l).Error
}

func GetAllUsers() ([]string, error) {
	var users []string
	err := D.DB.Table("users").Pluck("account", &users).Error
	return users, err
}

func GetTotalMoney(uname string) float64 {
	var sum float32
	if err := D.DB.Table("card_records").
		Where("account=?", uname).
		Select("SUM(trans_money) as sum").Scan(&sum).Error; err != nil {
		return -1
	}
	ret, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", sum), 32)
	return ret
}

func GetTotalRecord(uname string) int {
	var count int64
	D.DB.Table("card_records").
		Where("account=?", uname).
		Count(&count)
	return int(count)
}

func GetLocationTime(uname string) []OrderData {
	var location []OrderData
	
	D.DB.Table("card_records").
		Group("location").
		Select([]string{"location as name", "COUNT(location) as count"}).
		Where("account=?", uname).
		Order("count desc").
		Scan(&location)
	
	return location
}

func GetMethodTime(uname string) []OrderData {
	var method []OrderData
	
	D.DB.Table("card_records").
		Group("method").
		Select([]string{"method as name", "COUNT(method) as count"}).
		Where("account=?", uname).
		Order("count desc").
		Scan(&method)
	
	return method
}

func GetBooks(uname string) []string {
	var names []string
	D.DB.Table("books").
		Where("account=?", uname).
		Pluck("name", &names)
	return names
}

func (d *Database) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", viper.GetString("mysql.user"),
		viper.GetString("mysql.password"), viper.GetString("mysql.ip"),
		viper.GetString("mysql.port"), viper.GetString("mysql.database"))
}

func (d *Database) Init() {
	var err error
	fmt.Println(d.getDSN())
	d.DB, err = gorm.Open(mysql.Open(d.getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	if err := d.DB.AutoMigrate(&User{}, &Book{}, &CardRecord{}); err != nil {
		log.Fatal("表初始化失败:", err)
	}
}

func Start() {
	D.Init()
	DR.Init()
}
