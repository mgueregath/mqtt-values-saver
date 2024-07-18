package mariadb

import (
	"ValuesImporter/facade/environment"
	"ValuesImporter/persistence"
	"database/sql"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var mariadb *MariaDB

var lock = &sync.Mutex{}

type MariaDB struct {
	poolConnection persistence.PoolConnection
}

func NewMariaDB(ENVIRONMENT *environment.Environment) *persistence.PoolConnection {
	if mariadb == nil {
		lock.Lock()
		defer lock.Unlock()
		if mariadb == nil {
			poolConnection := connect(ENVIRONMENT)
			mariadb = &MariaDB{poolConnection: *poolConnection}
			return poolConnection
		}

	}
	return &mariadb.poolConnection
}

func GetConnection() *persistence.PoolConnection {
	return &mariadb.poolConnection
}

func connect(ENVIRONMENT *environment.Environment) *persistence.PoolConnection {

	dataSource := ENVIRONMENT.Database.Username
	dataSource += ":" + ENVIRONMENT.Database.Password
	dataSource += "@tcp(" + ENVIRONMENT.Database.Host
	dataSource += ":" + strconv.Itoa(ENVIRONMENT.Database.Port)
	dataSource += ")/" + ENVIRONMENT.Database.Database
	db, err := sql.Open(ENVIRONMENT.Database.Driver, dataSource)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(ENVIRONMENT.Database.MaxDatabaseConnections)
	db.SetMaxIdleConns(ENVIRONMENT.Database.MaxIdleConnections)

	poolConnection := persistence.NewPoolConnection(db)
	return poolConnection
}
