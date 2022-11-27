package crawler

import (
	"context"
	"log"
	"net/http"
	"sync"
	"year-end/crawler/task/card"
	"year-end/crawler/task/library"
	"year-end/model"
	"year-end/utils/errno"
)

type Engine func(ctx context.Context, uname, psd string, client *http.Client) error

func Crawler(ctx context.Context, uname, psd string, client *http.Client) error {
	var wg sync.WaitGroup
	var errCh = make(chan error)
	var done = make(chan struct{})
	childCtx, cancel := context.WithCancel(ctx)
	
	go func() {
		for {
			select {
			case err, ok := <-errCh:
				if ok && err != nil {
					if errno.Is(err, errno.ERR_CACHE_EXIST) {
						continue
					}
					log.Println("err:", err)
					cancel()
				}
			case <-done:
				return
			}
		}
	}()
	wg.Add(2)
	
	go worker(childCtx, errCh, uname, psd, library.GetHistoryBooks, &wg, client)
	go worker(childCtx, errCh, uname, psd, card.GetRecordsOfYear, &wg, client)
	
	wg.Wait()
	// 将已存在的数据存放的缓冲中
	redisSetting(childCtx, errCh, uname)
	done <- struct{}{}
	return nil
}

func redisSetting(ctx context.Context, errCh chan<- error, uname string) {
	
	// set redis process
	if err := model.DR.SetOneUserRecord(ctx, uname); err != nil {
		errCh <- err
		return
	}
}

func worker(ctx context.Context, errCh chan<- error, uname, psd string, engine Engine, wg *sync.WaitGroup, client *http.Client) {
	defer wg.Done()
	if err := engine(ctx, uname, psd, client); err != nil {
		errCh <- err
		return
	}
}
