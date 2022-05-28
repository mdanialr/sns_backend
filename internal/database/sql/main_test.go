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
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config using viper: %s", err)
	}
	var conf service.Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("failed to unmarshal config: %s", err)
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
