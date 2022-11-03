package crawler

import (
	"context"
	"log"
	"sync"
	"year-end/crawler/task/card"
	"year-end/crawler/task/library"
	"year-end/model"
)

type Engine func(ctx context.Context, uname, psd string)

func Crawler(ctx context.Context, uname, psd string) {
	var wg sync.WaitGroup
	var user = &model.UserModel{
		Password: psd,
		ID:       uname,
	}
	if err := user.AddUser(); err != nil {
		log.Fatal("添加用户失败", err)
		return
	}
	// 使用engine
	childCtx, _ := context.WithCancel(ctx)
	worker(childCtx, uname, psd, library.GetHistoryBooks, &wg)
	worker(childCtx, uname, psd, card.GetRecordsOfYear, &wg)
	
	// 将已存在的数据存放的缓冲中
	//redisSetting(ctx, &wg, uname)
	
	wg.Wait()
}

func redisSetting(ctx context.Context, wg *sync.WaitGroup, uname string) {
	wg.Add(1)
	// set redis process
	model.SetOneUserRecord(ctx, uname)
	defer wg.Done()
}

func worker(ctx context.Context, uname, psd string, engine Engine, wg *sync.WaitGroup) {
	wg.Add(1)
	engine(ctx, uname, psd)
	defer wg.Done()
}
