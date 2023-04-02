package postgresql

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGorm init new GORM using postgresql as the db and given viper config to
// get detailed information about the db such as host, port, username etc.
func NewGorm(v *viper.Viper, logger logger.Interface) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s",
		v.GetString("db.host"),
		v.GetInt("db.port"),
		v.GetString("db.name"),
		v.GetString("db.user"),
		v.GetString("db.pass"),
	)
	// open connection to db based on the constructed dsn connection string above
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger,
	})
}
