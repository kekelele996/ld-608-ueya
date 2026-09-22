package database

import (
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"groundTurn/src/config"
	"groundTurn/src/models"
)

// Open selects MySQL in containers and SQLite (SQLITE_PATH) for local tests.
func Open(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	if cfg.SQLitePath != "" {
		dialector = sqlite.Open(cfg.SQLitePath)
	} else {
		dialector = mysql.Open(cfg.DSN())
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if cfg.SQLitePath != "" {
		// Single writer in SQLite: wait for locks instead of failing mid-test.
		if err := db.Exec("PRAGMA busy_timeout = 5000;").Error; err != nil {
			return nil, err
		}
	} else {
		// The compose healthcheck can pass a moment before the application
		// account is fully ready; retry pinging for up to ~60s.
		if err := waitForMySQL(sqlDB, 30); err != nil {
			return nil, err
		}
	}

	if err := Migrate(db); err != nil {
		return nil, err
	}
	log.Printf("database ready (sqlite=%v)", cfg.SQLitePath != "")
	return db, nil
}

func waitForMySQL(sqlDB interface {
	Ping() error
}, attempts int) error {
	var lastErr error
	for i := 0; i < attempts; i++ {
		if lastErr = sqlDB.Ping(); lastErr == nil {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return lastErr
}

// Migrate keeps the schema aligned with the models; database/init.sql creates
// the database itself for MySQL 8.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.FlightTurnaround{},
		&models.GroundTask{},
		&models.GroundResource{},
		&models.ResourceBooking{},
		&models.DelayEvent{},
		&models.AuditLog{},
	)
}
