package database

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
	"github.com/mdanialr/sns_backend/internal/service"
)

var (
	testDB      *sql.DB
	testQueries *Queries
)

func TestMain(m *testing.M) {
	// read from ENV variables
	dbPort, _ := strconv.ParseUint(os.Getenv("DB_PORT"), 10, 16)
	dbConfig := service.DBPostgres{
		Name: os.Getenv("DB_NAME"),
		Host: os.Getenv("DB_HOST"),
		Port: uint(dbPort),
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASS"),
	}

	db, err := sql.Open(dbConfig.GetDriver(), dbConfig.GenerateConnectionString())
	if err != nil {
		log.Fatalf("failed to open connection to database: %s", err)
	}
	testDB = db
	defer testDB.Close()

	testQueries = New(testDB)

	os.Exit(m.Run())
}
