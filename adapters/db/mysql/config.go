package mysql

import "time"

// Config holds everything needed to open a MariaDB/MySQL connection pool.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string

	// Pool tuning
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (c Config) dsn() string {
	if c.Port == 0 {
		c.Port = 3306
	}
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=UTC&charset=utf8mb4",
		c.User, c.Password, c.Host, c.Port, c.DBName,
	)
}