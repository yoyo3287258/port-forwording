package web

import (
	"github.com/gin-gonic/gin"
	"port-forwording/web/routers"
)

func Run() {
	r := gin.Default()
	routers.InitRouters(r)
	err := r.Run(":8765")
	if err != nil {
		return
	}
}
