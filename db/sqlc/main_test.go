package db

import (
	"database/sql"
	"log"
	"os"
	"simplebank/util"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
  config, err := util.LoadConfig("../..")
	if err != nil {
    log.Fatal("unable to load config for viper: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to DB")
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
