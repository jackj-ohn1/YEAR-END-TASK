package card

import (
	"fmt"
	"github.com/jinzhu/gorm"
	jsoniter "github.com/json-iterator/go"
	"log"
	"net/http"
	"net/url"
	"time"
	"year-end/model"
	"year-end/tool/user"
)

const nonTOKEN = ""

// 开多个worker的时间开销与只开month个的时间开销无差
func workers(n int, data <-chan model.ConsumeModel, tx *gorm.DB) {
	for i := 0; i < n; i++ {
		go func() {
			for {
				select {
				case val, ok := <-data:
					if !ok {
						return
					}
					fmt.Println(val)
					if err := (&val).AddRecord(tx); err != nil {
						log.Fatal("record存放失败:", err)
						tx.Rollback()
						return
					}
				}
			}
		}()
	}
}

func tick() {
	count := 1
	for {
		time.Sleep(time.Second * 1)
		fmt.Println("等待时间", count, "s")
		count++
	}
}

func getOneCCNUToken(uname, psd string) (string, error) {
	var client *http.Client
	var err error
	
	if client, err = user.Login(uname, psd); err != nil {
		return nonTOKEN, err
	}
	
	req, _ := http.NewRequest("POST", hostUrl, nil)
	
	// client 存放了所有的token 和 cookie
	client.Do(req)
	
	// 自己按需求取出
	hostUrl, _ := url.Parse(hostUrl)
	return "Bearer " + client.Jar.Cookies(hostUrl)[2].Value, nil
}

func parseJson(uname string, body []byte, data chan<- model.ConsumeModel, done chan<- struct{}) {
	var m model.Record
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(body, &m)
	if err != nil {
		done <- struct{}{}
		return
	}
	
	// 多个receiver和sender无法得知终止时间
	for _, v := range m.Result.Rows {
		data <- v.ConvertToModel(uname)
	}
	defer func() {
		done <- struct{}{}
	}()
}
