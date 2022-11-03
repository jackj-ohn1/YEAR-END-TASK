package model

import (
	"context"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"log"
)

var DR *DatabaseRedis = &DatabaseRedis{Dr: new(redis.Client)}

type DatabaseRedis struct {
	Dr *redis.Client
}

// 记录
type WrapOrderData struct {
	Method   []OrderData `json:"method" redis:"method"`
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
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
}

// 一次将所有的加到缓存上，刚开启时使用
func SetAllUserRecord(ctx context.Context) {
	users, err := GetAllUsers()
	if err != nil {
		log.Fatal("获取用户失败", err)
		return
	}
	
	for _, val := range users {
		ret, err := DR.Dr.Exists(ctx, val).Result()
		if err != nil {
			log.Fatal("查询redis数据库错误:", err)
			return
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
		if err = DR.Dr.LPush(ctx, val, one).Err(); err != nil {
			log.Fatal("存放缓冲失败", err)
		}
	}
}

// 用户第一次登录使用就直接存数据库+存缓冲
func SetOneUserRecord(ctx context.Context, uname string) {
	ret, err := DR.Dr.Exists(ctx, uname).Result()
	if err != nil {
		log.Fatal("查询redis数据库错误:", err)
		return
	}
	if ret > 0 {
		return
	}
	var t = &WrapOrderData{
		Method:   GetMethodTime(uname),
		Location: GetLocationTime(uname),
		Money:    GetTotalMoney(uname),
		Times:    GetTotalRecord(uname),
		Books:    GetBooks(uname),
	}
	if err := DR.Dr.LPush(ctx, uname, t); err != nil {
		log.Fatal("存放缓冲失败:", err)
	}
}

func GetOneUserRecord(ctx context.Context, uname string) WrapOrderData {
	var t = WrapOrderData{} //
	if err := DR.Dr.LPop(ctx, uname).Scan(&t); err != nil {
		log.Fatal("取数据失败:", err)
	}
	return t
}
