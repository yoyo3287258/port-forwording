package common

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
	"port-forwording/internal/model"
)

func init() {
	database, err := initializeDatabase()
	if err != nil {
		panic(err)
	}
	DB = database
}

func initializeDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(".port_forwarding.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移数据库结构
	err = db.AutoMigrate(&model.PortForwarding{})
	if err != nil {
		return nil, err
	}
	log.Println("open sqlite success")
	return db, nil
}
