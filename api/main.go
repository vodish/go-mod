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

// func mysqlPing() {
// 	pingErr := dbx.Ping()
// 	if pingErr != nil {
// 		log.Fatal(pingErr)
// 	}

// 	fmt.Println("mysql ping ok.")
// }

func ginRouter() {

	router := gin.Default()
	router.GET("/test", hTest)
	router.GET("/users", hUserList)
	router.GET("/users/:email", hUserOne)
	router.GET("/usersn", hUsersParamName)
	// router.POST("/albums", postAlbums)

	router.Run(os.Getenv("SERVER"))
}

func hTest(c *gin.Context) {
	c.JSON(200, gin.H{"str": "строка", "int": 200})
}

// album represents data about a record album.
type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Start int    `json:"start"`
}

// список пользователей
func hUserList(c *gin.Context) {
	var users []User

	sql := `SELECT *  FROM user`

	err := dbx.Select(&users, sql)

	if err != nil {
		fmt.Println("err", err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"err": "data is not found from getUserList"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// один пользователь
func hUserOne(c *gin.Context) {
	var user User

	sql := `SELECT *  FROM user  WHERE email = ?`

	email := c.Param("email")
	err := dbx.Get(&user, sql, email)

	if err != nil {
		fmt.Println("err", err)
	}
	if user.Id == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"err": "data is not found from getUserOne"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func hUsersParamName(c *gin.Context) {
	var user User
	var users []User

	rows, err := dbx.NamedQuery(`
		SELECT
			*
		FROM
			user
		WHERE
			email = :email
		LIMIT
			3
		`,
		map[string]interface{}{
			"email": "site@taris.pro",
		})

	for rows.Next() {
		rows.StructScan(&user)
		users = append(users, user)
	}
	defer rows.Close()

	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, users)
}

func main() {

	setEnv()   // переменные окружения из файла .env
	setMysql() // подключиться к mysql
	// mysqlPing() // проверка подключения
	ginRouter() // запустить сервер
}
