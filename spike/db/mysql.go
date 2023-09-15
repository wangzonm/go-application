package db

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"xorm.io/xorm"
)

var MySQLClient *xorm.Engine

func Init() {
	var err error
	MySQLClient, err = xorm.NewEngine("mysql", "uangme_loan_pro:0rvy8v68L*4@tcp(147.139.185.188:3306)/life_service?charset=utf8")
	if err != nil {
		panic(err)
	} else {
		MySQLClient.SetMaxOpenConns(10)
		MySQLClient.SetMaxIdleConns(10)
		MySQLClient.SetMaxOpenConns(10)
		MySQLClient.ShowSQL(true)
		if err = MySQLClient.Ping(); err != nil {
			panic(err)
		}
		log.Printf("Connected to mysql: %+v", MySQLClient)
	}
}

type Pk struct {
	Id    int64 `json:"id" xorm:"id"`
	Count int64 `json:"count" xorm:"count"`
}
