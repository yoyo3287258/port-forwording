package model

type PortForwarding struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	ListenPort uint   `json:"listenPort" gorm:"unique" binding:"required"`
	TargetIp   string `json:"targetIp" binding:"required"`
	TargetPort uint   `json:"targetPort" binding:"required"`
}
