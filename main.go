package main

import (
	"database/sql"
	"log"

	"github.com/ahmedabzk/simple_bank/api"
	db "github.com/ahmedabzk/simple_bank/db/sqlc"
	"github.com/ahmedabzk/simple_bank/util"

	_ "github.com/lib/pq"
)



func main() {
	config, err := util.LoadConfig(".")
	if err != nil{
		log.Fatal("cannot read config files", err)
	}
	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to the db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
