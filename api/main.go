package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

	fmt.Println("mysql connected.")
}

func mysqlPing() {
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("mysql ping ok.")
}

func ginRouter() {
	router := gin.Default()
	router.GET("/albums", getAlboms)

	// router.POST("/albums", postAlbums)

	router.Run(os.Getenv("SERVER"))

}

// album represents data about a record album.
// type album struct {
// 	ID     string  `json:"id"`
// 	Title  string  `json:"title"`
// 	Artist string  `json:"artist"`
// 	Price  float64 `json:"price"`
// }

// // albums slice to seed record album data.
// var albums = []album{
// 	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// 	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// 	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }

func getAlboms(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, "alboms handler content")
}

func main() {

	setEnv() // переменные окружения из файла .env

	setMysql() // подключиться к mysql

	mysqlPing() // проверка подключения

	ginRouter() // запустить сервер
}
