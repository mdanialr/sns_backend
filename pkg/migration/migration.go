package migration

import (
	"log"
	"os"

	"github.com/mdanialr/sns_backend/internal/domain"
	conf "github.com/mdanialr/sns_backend/pkg/config"
	gormLogger "github.com/mdanialr/sns_backend/pkg/gorm"
	"github.com/mdanialr/sns_backend/pkg/migration/seeder"
	"github.com/mdanialr/sns_backend/pkg/postgresql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Run do run migration (creating all tables) and optionally run seeder if
// given param is true.
func Run(isSeeder bool) {
	db := initGorm()
	// get the sql db
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln("failed to get the DB instance from gorm:", err)
		return
	}
	defer sqlDB.Close()

	db.AutoMigrate(
		&domain.RegisteredOTP{},
	)
	if isSeeder {
		seeder.Run(db)
	}
}

func initGorm() *gorm.DB {
	// init viper config
	v, err := conf.InitConfigYml()
	if err != nil {
		log.Fatalln("failed to init config:", err)
	}
	// setup the logger for
	gormLog := gormLogger.New(os.Stdout, logger.Error)
	// init gorm using postgresql as the DB
	db, err := postgresql.NewGorm(v, gormLog)
	if err != nil {
		log.Fatalln("failed to init gorm with postgresql as the DB:", err)
		return nil
	}
	db.Config.DisableForeignKeyConstraintWhenMigrating = true

	return db
}
