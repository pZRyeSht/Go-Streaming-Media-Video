package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

var (
	dbConn *sql.DB
	err error
)

func init() {
	dbConn, err = sql.Open("mysql", "root:123456@/video_server?charset=utf8")
	if err!=nil {
		panic(err.Error())
	}
	fmt.Printf("dbConn +%v\n", dbConn)
}