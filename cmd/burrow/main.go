package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ManoloEsS/burrow_prototype/cli"
	"github.com/ManoloEsS/burrow_prototype/internal/config"
	"github.com/ManoloEsS/burrow_prototype/internal/engine"
	"github.com/ManoloEsS/burrow_prototype/internal/models"
	// _ "github.com/mattn/go-sqlite3"
)

const (
	DefaultPort = ":8080"
	// dbFile      = "./requests.db"
)

func main() {
	// //connect to database
	// db, err := sql.Open("sqlite3", dbFile)
	// if err != nil {
	// 	log.Fatalf("could not connect to database: %v", err)
	// }
	// defer db.Close()
	//

	cfg := &config.Config{
		DefaultPort: DefaultPort,
	}
	ctx := context.Background()

	newReq := &models.Request{}

	//get input
	cli.InputToReq(cfg, newReq)

	response, err := engine.GetRequest(ctx, newReq.Method, newReq.URL, nil)
	if err != nil {
		log.Fatalf("invalid arguments: %v", err)
	}
	fmt.Printf("request to %s with method %s returns %d\n", newReq.URL, newReq.Method, response.StatusCode)
}
