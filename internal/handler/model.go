package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"year-end/utils/errno"
)

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func DefaultError(c *gin.Context, code int, err error) {
	c.JSON(200, wrapResp(code, nil))
	log.Println("program info:", err)
}

func InternalServerError(c *gin.Context, err error) {
	log.Panicln("program panic:", err)
	c.JSON(500, wrapResp(errno.ERR_INTERNAL_SERVER_WRONG, nil))
}

func Success(c *gin.Context, data any) {
	c.JSON(200, wrapResp(errno.SUCCESS, data))
}

func Default(c *gin.Context, code int) {
	c.JSON(200, wrapResp(code, nil))
}

func wrapResp(code int, data any) Resp {
	return Resp{
		Code:    code,
		Message: errno.GetMessage(code),
		Data:    data,
	}
}
