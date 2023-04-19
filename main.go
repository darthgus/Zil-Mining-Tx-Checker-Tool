package main

import (
	"log"
	"os"

	//"strconv"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var (
	//lint:ignore U1000 ignore
	db *sqlx.DB
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalln("Error loading .env")
	}

	DSBlock := 28112

	TXBlock := (DSBlock + 1) * 100
	EndingTxBlock := TXBlock + 100

	connectionString := os.Getenv("DATABASE_URL")

	if connectionString == "" {
		log.Fatalln("Please pass the connection string using the -conn option")
	}

	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v\n", err)
	}

	// Grab the DS Block Miner Info
	InsertDSBlockMinerInfo(DSBlock, db)

	//Grab each block GetTxnBodiesForTxBlock
	for TXBlock < EndingTxBlock {
		InsertTXBlockTransactions(TXBlock, db)
		TXBlock++
	}
}
