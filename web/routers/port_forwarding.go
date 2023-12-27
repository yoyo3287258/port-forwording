package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"port-forwording/internal/common"
	"port-forwording/internal/model"
	"port-forwording/internal/network"
)

func InitPortForwardRouter(r *gin.Engine) {
	router := r.Group("portForward")

	router.GET("/list", listHandler)
	router.POST("/add", addHandler)
}

func addHandler(c *gin.Context) {
	var pfInfo model.PortForwarding
	err := c.ShouldBind(&pfInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	common.DB.Save(&pfInfo)
	go network.Create(pfInfo.ListenPort, pfInfo.TargetIp, pfInfo.TargetPort)
	c.JSON(200, gin.H{"data": pfInfo})
}

func listHandler(c *gin.Context) {
	var pfList []model.PortForwarding
	common.DB.Find(&pfList)
	c.JSON(200, gin.H{"data": pfList})
}
