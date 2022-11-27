package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"year-end/utils/errno"
	"year-end/utils/token"
)

func Parse() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.JSON(200, gin.H{
				"message": errno.GetMessage(errno.ERR_TOKEN_NOT_EXIST),
				"data":    nil,
				"code":    errno.ERR_TOKEN_NOT_EXIST,
			})
			c.Abort()
			return
		}
		if len(strings.Split(auth, ".")) != 3 {
			c.JSON(200, gin.H{
				"message": errno.GetMessage(errno.ERR_TOKEN_TYPE_WRONG),
				"data":    nil,
				"code":    errno.ERR_TOKEN_TYPE_WRONG,
			})
			c.Abort()
			return
		}
		claims, err := token.ParseToken(auth)
		if err != nil {
			c.JSON(200, gin.H{
				"message": errno.GetMessage(errno.ERR_TOKEN_WRONG),
				"data":    nil,
				"code":    errno.ERR_TOKEN_WRONG,
			})
			c.Abort()
			return
		}
		if claims.ExpiresAt < time.Now().Unix() {
			c.JSON(200, gin.H{
				"message": errno.GetMessage(errno.ERR_TOKEN_EXPIRED),
				"data":    nil,
				"code":    errno.ERR_TOKEN_EXPIRED,
			})
			c.Abort()
			return
		}
		
		c.Set("account", claims.Account)
		c.Next()
	}
}
