package service

import "fmt"

// DBPostgres holds data to connect to postgresql database.
type DBPostgres struct {
	User        string `mapstructure:"user"` // the user that would be used to connect to the database
	Pass        string `mapstructure:"pass"` // the password of the User that would be used to connect to the database
	Host        string `mapstructure:"host"` // the host where postgresql database is run
	Port        uint   `mapstructure:"port"` // the port where postgresql database is run
	Name        string `mapstructure:"name"` // the name of the database
	IsSSL       bool   `mapstructure:"ssl"`  // set ssl mode of this database connection
	isSSLString string // generate string version of IsSSL field
}

// GenerateConnectionString generate connection string that could be supplied to database/sql.Open pkg.
func (d *DBPostgres) GenerateConnectionString() string {
	if d.User == "" {
		d.User = "postgres"
	}
	if d.Pass == "" {
		d.Pass = "postgres"
	}
	if d.Host == "" {
		d.Host = "127.0.0.1"
	}
	if d.Port == 0 {
		d.Port = 5432
	}
	if d.Name == "" {
		d.Name = "postgres"
	}

	if !d.IsSSL {
		d.isSSLString = "disable"
	}
	if d.IsSSL {
		d.isSSLString = "enable"
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Pass,
		d.Host, d.Port,
		d.Name,
		d.isSSLString,
	)
}

// GetDriver return string of the driver name that could be supplied to database/sql.Open pkg.
func (d *DBPostgres) GetDriver() string {
	return "postgres"
}
