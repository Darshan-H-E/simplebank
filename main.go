package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"

	_ "github.com/lib/pq"
)

func main() {
  config, err := util.LoadConfig(".")
	if err != nil {
    log.Fatal("unable to load config for viper: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
    log.Fatal("cannot connect to DB: ", err)
	}

  store := db.NewStore(conn)
  server := api.NewServer(store)

  err = server.Start(config.ServerAddress)
	if err != nil {
    log.Fatal("unable to start server: ", err)
	}
}
