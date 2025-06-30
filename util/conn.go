package util

import (
	"database/sql"
	"flag"

	"github.com/rah-0/nabu"
	"xorm.io/xorm"
)

func GetConn() (*sql.DB, *xorm.Engine, error) {
	dbUser := flag.String("dbUser", "", "Required")
	dbPassword := flag.String("dbPassword", "", "Required")
	dbName := flag.String("dbName", "", "Required")
	dbIp := flag.String("dbIp", "", "Required")
	dbPort := flag.String("dbPort", "3306", "Required")
	flag.Parse()

	conn, err := sql.Open("mysql", *dbUser+":"+*dbPassword+"@tcp("+*dbIp+":"+*dbPort+")/"+*dbName)
	if err != nil {
		return nil, nil, nabu.FromError(err).Log()
	}

	connXorm, err := xorm.NewEngine("mysql", *dbUser+":"+*dbPassword+"@tcp("+*dbIp+":"+*dbPort+")/"+*dbName)
	if err != nil {
		return nil, nil, nabu.FromError(err).Log()
	}

	return conn, connXorm, nil
}
