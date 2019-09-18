package common

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
	defer DB.Close()
	return err
}
