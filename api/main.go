package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB
var err error

func setEnv() {
	if godotenv.Load("../.env") != nil {
		log.Fatal("Did not exist file ../.env")
	}
}

func setMysql() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DBHOST"),
		DBName:               os.Getenv("DBNAME"),
		AllowNativePasswords: true,
	}

	// Get a database handle.
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
}

func mysqlPing() {
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
}

func main() {

	setEnv() // переменные окружения из файла .env

	setMysql() // подключиться к mysql

	mysqlPing() // проверка подключения

	fmt.Println("Connected!")
}
