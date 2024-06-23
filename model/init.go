package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

var DB *gorm.DB

func Database(connstring string) {
	db, err := gorm.Open("mysql", connstring)
	if err != nil {
		panic("Mysql数据库连接错误")
	}
	fmt.Println("数据库连接成功")
	db.LogMode(true)
	if gin.Mode() == "release" {
		db.LogMode(false)
	}
	db.SingularTable(true)                       // 表名不加s
	db.DB().SetMaxIdleConns(20)                  // 设置连接池
	db.DB().SetMaxOpenConns(100)                 // 最大连接数
	db.DB().SetConnMaxLifetime(30 * time.Second) // 连接存活时间
	DB = db
	migrate()
}
