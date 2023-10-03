package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Bot struct {
	Id          int
	UserId      int
	Title       string
	Description string
	//Content     string
}

func (u *Bot) TableName() string {
	return "bot"
}

func read(db *gorm.DB, id string) *Bot {
	var bot []Bot
	db.Where("id = ?", id).First(&bot)
	return &bot[0]
}
func main() {
	dataSourceName := "root:liujun@tcp(39.106.49.188)/kob?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dataSourceName), nil)
	if err != nil {
		panic(err)
	}
	user := read(db, "1")
	if user != nil {
		fmt.Println(*user)
	} else {
		fmt.Println("数据库出错")
	}
}
