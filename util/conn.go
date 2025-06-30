package util

import (
	"flag"
)

func GetDsn() string {
	dbUser := flag.String("dbUser", "", "Required")
	dbPassword := flag.String("dbPassword", "", "Required")
	dbName := flag.String("dbName", "", "Required")
	dbIp := flag.String("dbIp", "", "Required")
	dbPort := flag.String("dbPort", "3306", "Required")
	flag.Parse()

	return *dbUser + ":" + *dbPassword + "@tcp(" + *dbIp + ":" + *dbPort + ")/" + *dbName
}
