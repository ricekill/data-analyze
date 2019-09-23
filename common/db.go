package common

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func OpenDb() (err error) {
	DB, err = sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
			Config.Db.User,
			Config.Db.Password,
			Config.Db.Host,
			Config.Db.Port,
			Config.Db.Name))
	if err != nil {
		return err
	}
	//defer DB.Close()

	DB.SetConnMaxLifetime(10*time.Second)  //最大连接周期，超过时间的连接就close
	DB.SetMaxOpenConns(100) //设置最大连接数
	DB.SetMaxIdleConns(10) //设置闲置连接数
	return err
}
