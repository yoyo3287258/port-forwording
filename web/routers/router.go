package routers

import (
	"github.com/gin-gonic/gin"
)

func InitRouters(r *gin.Engine) {
	r.Any("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "pong",
		})
	})
	InitPortForwardRouter(r)

}
