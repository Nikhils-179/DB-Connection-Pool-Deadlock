package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" 
)

var sqlDB *sql.DB

func OpenDBConnection() *sql.DB {
	if sqlDB == nil {
		host := "mysql"
		port := "3307"
		user := "root"
		password := "root"
		database := "test"
		var err error

		// Correctly define the variable sqlDB
		sqlDB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", user, password, host, port, database))
		if err != nil {
			panic(err)
		}

		sqlDB.SetMaxOpenConns(5)  // Set max open connections at any time = 10
		sqlDB.SetMaxIdleConns(0)   // Set max idle (open, not in use) connections at any time = 5
	}

	return sqlDB
}


// hey -n 100 -c 20 -m GET -H "x-user-id: 1" http://localhost:4000/list-following
// -n  : number of Total Requests 
// -c : number of Concurrent Clients  
//increase the number of requests and increase number of 