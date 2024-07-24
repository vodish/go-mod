package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var dbx *sqlx.DB

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

	dbx, err = sqlx.Connect("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalln("dbx connect:", err)
	}

	fmt.Println("mysql sqlx connected.")
}

func mysqlPing() {
	pingErr := dbx.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("mysql ping ok.")
}

func ginRouter() {

	router := gin.Default()
	router.GET("/test", getTest)
	router.GET("/usersx", getUsersX)
	router.GET("/users", getUsers)
	// router.POST("/albums", postAlbums)

	router.Run(os.Getenv("SERVER"))
}

func getTest(c *gin.Context) {
	c.JSON(200, gin.H{"str": "строка", "int": 200})
}

// album represents data about a record album.
type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Start int    `json:"start"`
}

func getUsersX(c *gin.Context) {
	var users []User

	sql := "SELECT *  FROM user  LIMIT 1"
	err := dbx.Select(&users, sql)
	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println("users", users)

	// var users []User

	// for rows.Next() {
	// 	var user User
	// 	if err := rows.Scan(&user.id, &user.email, &user.start); err != nil {
	// 		c.JSON(400, gin.H{"error": err})
	// 	}
	// 	users = append(users, user)
	// }

	// if err := rows.Err(); err != nil {
	// 	c.JSON(500, gin.H{"error": err})
	// }

	c.JSON(http.StatusOK, users)
}

func getUsers(c *gin.Context) {
	var err error

	rows, err := dbx.NamedQuery(`
		SELECT
			*
		FROM
			user
		WHERE
			email = :email
		`,
		map[string]interface{}{
			"email": "site@taris.pro",
		})

	defer rows.Close()

	if err != nil {
		c.JSON(400, gin.H{"error": err})
	}

	var user User
	err = dbx.Get(&user, "SELECT *  FROM user  WHERE email = ?", "site@taris.pro")
	if err != nil {
		c.JSON(400, gin.H{"error rows next": err})
	}

	c.JSON(http.StatusOK, user)
}

func main() {

	setEnv()    // переменные окружения из файла .env
	setMysql()  // подключиться к mysql
	mysqlPing() // проверка подключения
	ginRouter() // запустить сервер
}
