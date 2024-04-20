package database

import (
	"fmt"
	"go-server/pkg/shared/logging"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
)

// NewDB creates a new database connection
func NewDB(
	logger *logrus.Logger,
) (*gorm.DB, error) {
	var dbConn *gorm.DB
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	dbConn, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}

	// create the database if it does not exist
	err = dbConn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", os.Getenv("DB_NAME"))).Error
	if err != nil {
		return nil, err
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	// create a new database connection with the logger
	dbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logging.NewGormLogger(logging.GormLogger{
			Logger:                    logger.WithField("service", "database"),
			LogLevel:                  gormLog.Info,
			IgnoreRecordNotFoundError: false,
			SlowThreshold:             200 * time.Millisecond,
			FileWithLineNumField:      "caller",
		}),
	})
	if err != nil {
		return nil, err
	}

	// ping the database to ensure the connection is valid
	err = Ping(dbConn)
	return dbConn, err
}

// CloseDB closes the database connection pool
func CloseDB(
	logger *logrus.Logger, // logger instance
	db *gorm.DB, // database instance
) {
	// get the *sql.DB from the gorm.DB
	myDB, err := db.DB()
	if err != nil {
		logger.Errorf("Error while returning *sql.DB: %v", err)
	}

	// log that we are closing the connection pool
	logger.Info("Closing the DB connection pool")

	// close the connection pool
	if err := myDB.Close(); err != nil {
		logger.Errorf("Error while closing the master DB connection pool: %v", err)
	}
}

// Ping checks if the database connection is alive.
func Ping(db *gorm.DB) error {
	// get the *sql.DB from the gorm.DB
	myDB, err := db.DB()
	if err != nil {
		return err
	}

	// ping the database
	return myDB.Ping()
}
