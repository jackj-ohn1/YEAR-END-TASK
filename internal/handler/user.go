package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"year-end/crawler"
	"year-end/model"
	"year-end/utils/errno"
	"year-end/utils/token"
	"year-end/utils/user"
)

type login struct {
	Account  string `form:"account"`
	Password string `form:"password"`
}

// 登陆成功就插入数据库，插入之前得先判断数据库里是否存在这条数据
func LoginAPI(c *gin.Context) {
	var one login
	if err := c.ShouldBindWith(&one, binding.Form); err != nil {
		DefaultError(c, errno.ERR_JSON_BIND_WRONG, err)
		return
	}
	client, err := user.Login(one.Account, one.Password)
	if err != nil {
		DefaultError(c, errno.ERR_ACCOUNT_PASSWORD_WRONG, err)
		return
	}
	
	var t = &model.User{Account: one.Account, Password: one.Password}
	if err := t.AddUser(); err != nil && !errno.Is(err, errno.ERR_USER_EXIST) {
		InternalServerError(c, err)
		return
	}
	
	// 缓存不存在就开启爬虫
	_, err = model.DR.GetOneUserRecord(c.Request.Context(), one.Account)
	if err != nil && errno.Is(err, errno.ERR_CACHE_NOT_EXIST) {
		go func() {
			if err := crawler.Crawler(context.Background(), one.Account, one.Password, client); err != nil {
				return
			}
		}()
	} else if err != nil {
		InternalServerError(c, err)
		return
	}
	
	// 返回token
	auth, err := token.GenerateToken(one.Account)
	if err != nil {
		InternalServerError(c, err)
		return
	}
	
	Success(c, auth)
}

func GetData(c *gin.Context) {
	// token与缓冲同步过期，所以在获取数据这里不需要额外再开启爬虫
	account := c.MustGet("account").(string)
	data, err := model.DR.GetOneUserRecord(c.Request.Context(), account)
	if err != nil {
		if !errno.Is(err, errno.ERR_CACHE_NOT_EXIST) {
			Default(c, errno.ERR_CACHE_NOT_EXIST)
			return
		}
		InternalServerError(c, err)
		return
	}
	Success(c, data)
}
