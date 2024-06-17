package database

import "fmt"

type dbinfo struct {
	name     string
	host     string
	port     string
	username string
	password string
}

func setupDB(name, host, port, username, password string) *dbinfo {
	return &dbinfo{
		name,
		host,
		port,
		username,
		password,
	}
}

func (db dbinfo) getStringPath() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Montreal", db.host, db.username, db.password, db.name, db.port)
}

func (db dbinfo) getURI() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.username, db.password, db.host, db.port, db.name)
}
