package main

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/raoptimus/gserv/config"
)

type (
	context struct {
		DB *sql.DB
	}
)

var Context *context

func initMysqlDb(urlRaw string) (*sql.DB, error) {
	db, err := sql.Open("mysql", urlRaw)
	if err != nil {
		return nil, errors.New("Couldn't connect to mysqldb (" + urlRaw + "): " + err.Error())
	}
	return db, nil
}

func Init() {
	if Context != nil {
		return
	}
	Context = &context{}
	if err := connect(); err != nil {
		log.Fatalln(err)
	}
	config.OnAfterLoad("data.Context.reconnect", reconnect)
	log.Println("data.Context initied")
}

func reconnect() {
	err := connect()
	if err != nil {
		log.Println(err)
	}
}

func connect() error {
	log.Println("data.Context db connection...")

	dbUrl := config.String("MySqlServer", "")
	if dbUrl != "" {
		db, err := initMysqlDb(dbUrl)
		if err != nil {
			return err
		}
		Context.DB = db
	}
	return nil
}
