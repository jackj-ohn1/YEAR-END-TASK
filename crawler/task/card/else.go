package card

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"year-end/model"
)

const nonTOKEN = ""

// 开多个worker的时间开销与只开month个的时间开销无差
func workers(n int, data <-chan model.CardRecord, tx *gorm.DB) {
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
						return
					}
				}
			}
		}()
	}
}

func getOneCCNUToken(uname, psd string, client *http.Client) (string, error) {
	//var client *http.Client
	//var err error
	//
	//if client, err = user.Login(uname, psd); err != nil {
	//	return nonTOKEN, err
	//}
	
	req, _ := http.NewRequest("POST", hostUrl, nil)
	
	// client 存放了所有的token 和 cookie
	if _, err := client.Do(req); err != nil {
		return nonTOKEN, err
	}
	
	// 自己按需求取出
	hostUrl, _ := url.Parse(hostUrl)
	return "Bearer " + client.Jar.Cookies(hostUrl)[2].Value, nil
}

func parseJson(uname string, body []byte, data chan<- model.CardRecord) error {
	var m model.Record
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	
	// 多个receiver和sender无法得知终止时间
	for _, v := range m.Result.Rows {
		data <- v.ConvertToModel(uname)
	}
	return nil
}
