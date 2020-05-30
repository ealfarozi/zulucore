package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"

	//get mysql driver
	_ "github.com/go-sql-driver/mysql"
)

//ViperEnvVariable is func to get .env file
func ViperEnvVariable(key string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value, ok := viper.Get(key).(string)

	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}

//InitializeMySQL is the func to open the connection for MySQL
func InitializeMySQL() *sql.DB {
	dBConnection, err := sql.Open("mysql", ViperEnvVariable("DB_URL"))

	if err != nil {
		fmt.Println("Connection Failed!!")
	}
	err = dBConnection.Ping()
	if err != nil {
		fmt.Println("Ping Failed!!")
	}

	dBConnection.SetMaxOpenConns(10)
	dBConnection.SetMaxIdleConns(5)
	dBConnection.SetConnMaxLifetime(time.Second * 10)
	return dBConnection
}

//CloseRows is func to close the rows result from SELECT query
func CloseRows(rows *sql.Rows) {
	if rows != nil {
		rows.Close()
	}
}

//CloseStmt is func to prepare the statement
func CloseStmt(stmt *sql.Stmt) {
	if stmt != nil {
		stmt.Close()
	}
}
