package model

type port_forwarding struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	ListenPort uint   `json:"listenPort" gorm:"unique" binding:"require"`
	TargetIp   string `json:"targetIp" binding:"require"`
	TargetPort uint   `json:"targetPort" binding:"require"`
}
