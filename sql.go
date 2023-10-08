package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const CREATE_DB_STATMENT string = "CREATE TABLE IF NOT EXISTS short_code (code text not null primary key, real_url text not null);"
const SELECT_REAL_URL_STATMENT string = "SELECT real_url FROM short_code WHERE code = ?;"
const INSERT_URL_STATMENT string = "INSERT INTO short_code(code, real_url) VALUES (?, ?);"

func initSQL() {
	db, err := sql.Open("sqlite3", "./shortCode.db")
	if err != nil {
		log.Fatal("Cant connect to db")
	}
	defer db.Close()
	db.Exec(CREATE_DB_STATMENT)

}

func queryCode(code string) (realURL string, err error) {
	db, err := sql.Open("sqlite3", "./shortCode.db")
	if err != nil {
		log.Fatal("Cant connect to db")
	}
	defer db.Close()
	result, err := db.Query(SELECT_REAL_URL_STATMENT, code)
	if err != nil {
		return "", err
	}
	defer result.Close()
	if result.Next() {
		err = result.Scan(&realURL)
	} else {
		return "404", errors.New("not found")
	}
	if err != nil {
		return "", err
	}
	return realURL, nil
}

func insertURL(code string, url string) (err error) {
	db, err := sql.Open("sqlite3", "./shortCode.db")
	if err != nil {
		log.Fatal("Cant connect to db")
	}
	defer db.Close()
	_, err = db.Exec(INSERT_URL_STATMENT, code, url)
	return err
}
