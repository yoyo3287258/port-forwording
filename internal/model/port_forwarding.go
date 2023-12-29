package model

type PortForwarding struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	ListenPort uint   `json:"listenPort" form:"listenPort" gorm:"unique" binding:"required"`
	TargetIp   string `json:"targetIp" form:"targetIp" binding:"required"`
	TargetPort uint   `json:"targetPort" form:"targetPort" binding:"required"`
	Remark     string `json:"remark" form:"remark" binding:"required"`
}
