package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/mdanialr/sns_backend/internal/service"
	"github.com/spf13/viper"
)

var (
	testDB      *sql.DB
	testQueries *Queries
)

func TestMain(m *testing.M) {
	viper.AddConfigPath("../../../")
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	var conf service.Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("failed to unmarshal config: %s", err)
	}

	viper.AutomaticEnv()

	// try to read from ENV variables
	dbDriver := os.Getenv("db_driver")
	dbSrc := os.Getenv("db_source")

	if dbDriver != "" && dbSrc != "" {
		conf.DBDriver = dbDriver
		conf.DBSource = dbSrc
	}

	db, err := sql.Open(conf.DBDriver, conf.DBSource)
	if err != nil {
		log.Fatalf("failed to open connection to database: %s", err)
	}
	testDB = db
	defer testDB.Close()

	testQueries = New(testDB)

	os.Exit(m.Run())
}
