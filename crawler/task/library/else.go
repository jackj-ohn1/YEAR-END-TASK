package library

import (
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func getHistoryBooks(cookie string, data chan<- string) {
	
	req, _ := http.NewRequest("GET", libraryUrl, nil)
	req.Header.Add("Cookie", cookie)
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("请求失败:", err)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	regexpMatch(body, data)
}

func regexpMatch(body []byte, data chan<- string) {
	//<td bgcolor="#FFFFFF" class="whitetext" width="25%"><a class="blue" href=".*?">(.*?)</a></td>
	regText := `<td bgcolor="#FFFFFF" class="whitetext" width="25%"><a class="blue" href=".*?">(.*?)</a></td>`
	reg := regexp.MustCompile(regText)
	ret := reg.FindAllSubmatch(body, -1)
	for k := range ret {
		data <- html.UnescapeString(string(ret[k][1]))
	}
	defer func() {
		close(data)
	}()
}

func getLibraryCookie(client *http.Client) (string, error) {
	//var client *http.Client
	//var err error
	//
	//if client, err = user.Login(uname, psd); err != nil {
	//	return noneCOOKIE, err
	//}
	
	req, _ := http.NewRequest("GET", libraryUrl, nil)
	
	client.Do(req)
	parseUrl, _ := url.Parse(hostUrl)
	
	return client.Jar.Cookies(parseUrl)[0].String(), nil
}
