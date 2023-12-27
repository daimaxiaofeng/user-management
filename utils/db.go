package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const CONFIG_FILE_PATH = "./config/db.json"

type Config struct {
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

func (c *Config) readFromFile(filepath string) error {
	configData, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(configData, &c)
}

var DB *sql.DB

func init() {
	var config Config
	if err := config.readFromFile(CONFIG_FILE_PATH); err != nil {
		panic(err)
	}
	if err := ConnectDB(config); err != nil {
		panic(err)
	}
}

func ConnectDB(c Config) error {
	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		c.Username, c.Password, c.Host, c.Port, c.Database)
	var err error
	DB, err = sql.Open("mysql", dbConnectionString)
	if err != nil {
		return err
	}
	return DB.Ping()
}
