package dao

import (
	"fmt"
	"test/configs"
	"test/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	cfg := configs.Cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Homework{},
		&models.Submission{},
		&models.Exam{},
	)
	if err != nil {
		return err
	}

	DB = db
	return nil
}
