package common

import (
	"database/sql"
	"log"
)

var DB *sql.DB
var Log *Logger
var Config *ServerConfig

func LoadConfig() error {
	Config = &ServerConfig{}
	return Config.load()
}
func SetupLogger() error {
	var err error
	Log, err = NewLogger(Config.Log.LogFile, Config.Log.TraceLevel)
	return err
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
