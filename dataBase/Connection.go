package dataBase

import (
	"database/sql"
	"fmt"
	"helloWorld/rule"
)

func CreateConnection() (*sql.DB, error) {
	dataBase, err := sql.Open(rule.SQL_DRIVER, rule.SQL_PATH)

	if err != nil {
		fmt.Println("DataBase create connection error : ", err)
	} else {
		fmt.Println("The database connection is successfully created")
	}

	return dataBase, err
}

func CloseConnection(dataBase *sql.DB) {
	err := dataBase.Close()
	if err != nil {
		fmt.Println("DataBase release connection error : ", err)
		return
	}
}
