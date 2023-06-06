package common

import (
	"fmt"
	"gin/model"
	"net/url"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	//可以用navicat或datagram等数据库操作软件，利用下面的信息登录数据库查看数据
	host := "gz-cynosdbmysql-grp-04b3z61j.sql.tencentcdb.com"
	port := "20464"
	database := "gin" //使用了mysql中的gin数据库，不要改其他的数据库！！
	username := "root"
	password := "Wobujidemimale123"
	charset := "utf8"
	loc := "Asia/Shanghai"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error to Db connection, err: " + err.Error())
	}
	db.AutoMigrate(&model.User{}) // 此处创建了model文件夹下的user实体类，仅作参考
	db.AutoMigrate(&model.Goods{})
	db.AutoMigrate(&model.Category{})
	db.AutoMigrate(&model.Banner{})
	db.AutoMigrate(&model.Picture{})
	db.AutoMigrate(&model.Chat{})
	db.AutoMigrate(&model.ChatList{})
<<<<<<< HEAD
	db.AutoMigrate(&model.Cart{})
=======
	db.AutoMigrate(&model.Order{})
	db.AutoMigrate(&model.Image{})
>>>>>>> 2f7999f1a791b3d5437a6c64f80f409d0cf26d0c

	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
