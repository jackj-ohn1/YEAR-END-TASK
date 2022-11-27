package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"log"
	"time"
	"year-end/utils/errno"
)

var DR *DatabaseRedis = &DatabaseRedis{Dr: new(redis.Client)}

type DatabaseRedis struct {
	Dr *redis.Client
}

// 记录
type WrapOrderData struct {
	// 支付方式
	Method []OrderData `json:"method" redis:"method"`
	// 支付地点
	Location []OrderData `json:"location" redis:"location"`
	Money    float64     `json:"money" redis:"money"`
	Times    int         `json:"times" redis:"times"`
	Books    []string    `json:"books"`
}

func (o *WrapOrderData) MarshalBinary() (data []byte, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(o)
}

func (o *WrapOrderData) UnmarshalBinary(data []byte) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, o)
}

func (dr *DatabaseRedis) Init() {
	dr.Dr = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", viper.GetString("redis.ip"), viper.GetString("redis.port")),
		DB:   viper.GetInt("redis.db"),
	})
	if err := dr.Dr.Ping(context.Background()).Err(); err != nil {
		log.Fatal("缓存设置失败:", err)
	}
}

// 一次将所有用户的缓存加上，刚开启时使用，还要一个全校统计的数据
func (dr *DatabaseRedis) SetAllUserRecord(ctx context.Context) error {
	users, err := GetAllUsers()
	if err != nil {
		return err
	}
	
	for _, val := range users {
		ret, err := dr.Dr.Exists(ctx, val).Result()
		if err != nil {
			return err
		}
		if ret > 0 {
			continue
		}
		
		var one = &WrapOrderData{
			Method:   GetMethodTime(val),
			Location: GetLocationTime(val),
			Money:    GetTotalMoney(val),
			Times:    GetTotalRecord(val),
			Books:    GetBooks(val),
		}
		if err = dr.Dr.Set(ctx, val, one, time.Hour*2).Err(); err != nil {
			return err
		}
	}
	return nil
}

// 用户第一次登录使用就直接存数据库+存缓冲
func (dr *DatabaseRedis) SetOneUserRecord(ctx context.Context, uname string) error {
	ret, err := dr.Dr.Exists(ctx, uname).Result()
	if err != nil {
		return err
	}
	if ret > 0 {
		return errors.New(errno.GetMessage(errno.ERR_CACHE_EXIST))
	}
	var t = &WrapOrderData{
		Method:   GetMethodTime(uname),
		Location: GetLocationTime(uname),
		Money:    GetTotalMoney(uname),
		Times:    GetTotalRecord(uname),
		Books:    GetBooks(uname),
	}
	if err := dr.Dr.Set(ctx, uname, t, time.Hour*2).Err(); err != nil {
		return err
	}
	
	return nil
}

func (dr *DatabaseRedis) GetOneUserRecord(ctx context.Context, uname string) (*WrapOrderData, error) {
	ret, err := dr.Dr.Exists(ctx, uname).Result()
	if err != nil {
		return nil, err
	}
	
	if ret <= 0 {
		return nil, errors.New(errno.GetMessage(errno.ERR_CACHE_NOT_EXIST))
	}
	
	var t WrapOrderData //
	if err := dr.Dr.Get(ctx, uname).Scan(&t); err != nil {
		return nil, err
	}
	return &t, nil
}
