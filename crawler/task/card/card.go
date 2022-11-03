package card

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"year-end/model"
)

const cardUrl = "http://one.ccnu.edu.cn/ecard_portal/query_trans"
const hostUrl = "http://one.ccnu.edu.cn"

func GetRecordsOfYear(ctx context.Context, uname, psd string) {
	token, err := getOneCCNUToken(uname, psd)
	if err != nil || token == nonTOKEN {
		log.Println(err)
		return
	}
	fmt.Println("获取token成功")
	go tick()
	
	data, done := make(chan model.ConsumeModel), make(chan struct{})
	nowYear, nowMonth := time.Now().Year(), time.Now().Month()
	tx := model.D.DB.Begin()
	workers(5, data, tx)
	for month := 1; month <= int(nowMonth); month++ {
		start, end := fmt.Sprintf("%d-0%d-0%d", nowYear, month, 1), fmt.Sprintf("%d-0%d-0%d", nowYear, month+1, 1)
		go getConsumptionRecords(uname, token, start, end, data, done)
	}
	fmt.Println("开始爬取数据...")
	i := 0
	for {
		select {
		case <-done:
			i++
			if i == int(nowMonth) {
				tx.Commit()
				close(data)
				return
			}
		case <-ctx.Done():
			tx.Rollback()
			return
		}
	}
}

func getConsumptionRecords(uname, token, start, end string, data chan model.ConsumeModel, done chan struct{}) {
	vals := url.Values{}
	// 一个月的个数
	vals.Set("limit", "1000")
	vals.Set("page", "1")
	vals.Set("start", start)
	vals.Set("end", end)
	vals.Set("tranType", "")
	
	req, err := http.NewRequest("POST", cardUrl, strings.NewReader(vals.Encode()))
	if err != nil {
		log.Println(err)
		return
	}
	
	req.Header.Add("Authorization", token)
	req.Header.Add("Referer", "http://one.ccnu.edu.cn/index")
	req.Header.Add("Host", "one.ccnu.edu.cn")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	
	parseJson(uname, body, data, done)
}
