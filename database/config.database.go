package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDb() error {
	var err error

	fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	fmt.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))
	fmt.Println("DB_PORT:", os.Getenv("DB_PORT"))

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable TimeZone=Asia/Jakarta"
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("error while connecting to database:", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("error while pinging to database:", err)
	}
	log.Println("successfully connected to database")
	return nil
}
