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
	router.DELETE("/del", delHandler)
}

func delHandler(c *gin.Context) {
	pfInfo := model.PortForwarding{}
	err := c.BindJSON(&pfInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if pfInfo.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	common.DB.Delete(&pfInfo)
	c.JSON(200, gin.H{"data": pfInfo, "code": 200})

}

func addHandler(c *gin.Context) {
	pfInfo := model.PortForwarding{}
	err := c.ShouldBind(&pfInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	common.DB.Save(&pfInfo)
	go network.Create(pfInfo.ListenPort, pfInfo.TargetIp, pfInfo.TargetPort)
	c.JSON(200, gin.H{"data": pfInfo, "code": 200})
}

func listHandler(c *gin.Context) {
	var pfList []model.PortForwarding
	common.DB.Find(&pfList)
	c.JSON(200, gin.H{"code": 0, "data": pfList, "count": len(pfList)})
}
