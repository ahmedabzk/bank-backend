package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/ahmedabzk/simple_bank/util"
	_ "github.com/lib/pq"
)


var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
		config,err := util.LoadConfig("../../")
	if err != nil{
		log.Fatal("cannot read config files:", err)
	}
	testDb, err = sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to the db", err)
	}
	testQueries = New(testDb)

	os.Exit(m.Run())
}
