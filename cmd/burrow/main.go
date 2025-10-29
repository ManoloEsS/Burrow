package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ManoloEsS/Burrow/config"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//initialize database
	db, err := sql.Open("sqlite3", "../../sql/mydb.sqlite")
	if err != nil {
		log.Fatal("Couldn't open database file")
	}
	defer db.Close()

	//intialize config
	config := config.ClientConfig{
		Db: db,
	}

	//initialize tui

	//create client
	client := &http.Client{}
}
