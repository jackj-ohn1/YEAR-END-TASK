package user

import (
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const loginUrl = "https://account.ccnu.edu.cn/cas/login"
const loginHost = "https://account.ccnu.edu.cn"

func Login(uname, psd string) (*http.Client, error) {
	var login bool
	jar, _ := cookiejar.New(&cookiejar.Options{})
	client := &http.Client{
		Timeout: 5 * time.Second,
		Jar:     jar,
	}
	
	htmlLogin, err := http.Get(loginUrl)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(htmlLogin.Body)
	if err != nil {
		return nil, err
	}
	
	if login, err = newRequest(client, body, uname, psd); err != nil {
		return nil, err
	} else if !login {
		return nil, errors.New("登录失败")
	}
	
	return client, nil
}

func newRequest(client *http.Client, body []byte, uname, psd string) (bool, error) {
	vals := getFormDate(body, uname, psd)
	
	req, err := http.NewRequest("POST", loginHost+getRegexpResult(`action="(.*?)" method="post"`, body), strings.NewReader(vals.Encode()))
	if err != nil {
		return false, nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	
	// 失败后跳转也是200
	b, _ := io.ReadAll(resp.Body)
	if !loginIn(b) {
		return false, nil
	}
	return true, nil
}

func getFormDate(body []byte, uname, psd string) url.Values {
	vals := url.Values{}
	vals.Set("_eventId", getRegexpResult(`name="_eventId" value="(.*?)"`, body))
	vals.Set("execution", getRegexpResult(`name="execution" value="(.*?)"`, body))
	vals.Set("lt", getRegexpResult(`name="lt" value="(.*?)"`, body))
	vals.Set("password", psd)
	vals.Set("username", uname)
	return vals
}

func getRegexpResult(rgx string, body []byte) string {
	rgxPattern := regexp.MustCompile(rgx)
	return string(rgxPattern.FindAllSubmatch(body, 1)[0][1])
}

func loginIn(body []byte) bool {
	rgxPattern := regexp.MustCompile("登录成功")
	return rgxPattern.Match(body)
}
