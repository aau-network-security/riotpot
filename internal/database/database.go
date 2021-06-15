package database

import (
	"fmt"

	"github.com/riotpot/internal/configuration"
	"github.com/riotpot/tools/errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Struct that defines the database connection.
// the database connection uses a pool so multiple
// services can push information to the database simultaneously
// through riotpot.
type Database struct {
	// Database configuration
	Config configuration.ConfigDatabase

	// Pointer to the db connection
	conn *gorm.DB
}

// Method to create a connection pool to the database
func (db *Database) connect() (*gorm.DB, error) {
	// build the connection url as a string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		db.Config.Host,
		db.Config.Username,
		db.Config.Password,
		db.Config.Identity.Name,
		db.Config.Port,
	)

	// create the pool connection to the database using its configuration.
	// It uses pgx as a default driver.
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	errors.Raise(err)

	db.conn = conn

	// returns the connection pool as a pointer
	return db.conn, err
}

// Get the connection to the database or create a new one.
func (db *Database) Connection() *gorm.DB {
	// check if there is any connection at all
	if db.conn != nil {
		return db.conn
	} else {
		conn, err := db.connect()
		errors.Raise(err)
		return conn
	}
}
