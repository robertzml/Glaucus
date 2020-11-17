package tdengine

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/taosdata/driver-go/taosSql"
)

var (
	DRIVER_NAME    = "taosSql"
	user           = "root"
	password       = "taosdata"
	host           = "47.111.23.211"
	port           = 6030
	dbName         = "Molan"
	dataSourceName = fmt.Sprintf("%s:%s@/tcp(%s:%d)/%s?interpolateParams=true", user, password, host, port, dbName)
	total          = 0
	lock           sync.Mutex
	nThreads       = 10
	nRequests      = 10
	profile        = "CPU.profile"
)

type Repository struct {
	db *sql.DB
}

func open() {
	db, err := sql.Open(DRIVER_NAME, dataSourceName)
	if err != nil {
		log.Fatal(err.Error())
	}

	rows, err := db.Exec("show tables")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(rows)
}