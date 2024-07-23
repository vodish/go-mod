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

func setEnv() {
	if godotenv.Load("../.env") != nil {
		log.Fatal("Did not exist file ../.env")
	}
}

func setMysql() {
	var err error

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
	router.GET("/test", getTest)
	router.GET("/users", getUsers)
	// router.POST("/albums", postAlbums)

	router.Run(os.Getenv("SERVER"))
}

// album represents data about a record album.
type User struct {
	id    string // `json:"id"`
	email string // `json:"email"`
	start int    //`json:"start"`
}

// // albums slice to seed record album data.
// var users = []user{
// 	{Id: "1", Email: "test@test", Start: 11},
// 	{Id: "2", Email: "2@test", Start: 22},
// }

func getTest(c *gin.Context) {
	c.JSON(200, gin.H{"str": "строка", "int": 200})
}

func getUsers(c *gin.Context) {
	email := "site@taris.pro"
	sql := `
		SELECT
			*
		FROM
			users
		WHERE
			email = ?
	`
	rows, err := db.Query(sql, email)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.id, &user.email, &user.start); err != nil {
			c.JSON(400, gin.H{"error": err})
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, users)
}

func main() {

	setEnv()    // переменные окружения из файла .env
	setMysql()  // подключиться к mysql
	mysqlPing() // проверка подключения
	ginRouter() // запустить сервер
}
