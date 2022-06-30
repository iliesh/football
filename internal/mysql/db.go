package dbx

import (
	"database/sql"
	"sync"
	"time"

	log "github.com/iliesh/go-templates/logger"

	_ "github.com/go-sql-driver/mysql"
)

func Open(host, port, name, user, pass string) (*sql.DB, error) {
	log.Debug("Connecting to the Database: Host: <%s>, Port: <%s>, Name: <%s>, User: <%s>, Password: <*****>", host, port, name, user)

	var dbOnce sync.Once
	var db *sql.DB
	var err error

	dbOnce.Do(func() { db, err = dbOpen(host, port, name, user, pass) })
	if err != nil {
		log.Trace("Error accessing the DB, return")
		return db, err
	}

	return db, nil
}

func dbOpen(host, port, name, user, pass string) (*sql.DB, error) {
	db, err := sql.Open("mysql", user+":"+pass+"@tcp("+host+":"+port+")/"+name)
	if err != nil {
		log.Error("Invalid DB config: %v", err)
		return nil, err
	}

	log.Trace("Connection was established, check if the DB is accessible")
	if err = db.Ping(); err != nil {
		log.Error("DB unreachable: %v", err)
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}
