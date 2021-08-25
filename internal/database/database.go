package database

import (
	// "fmt"
	"time"
	"context"

	"github.com/riotpot/internal/configuration"
	"github.com/riotpot/tools/errors"

	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Struct that defines the database connection.
// the database connection uses a pool so multiple
// services can push information to the database simultaneously
// through riotpot.
type Database struct {
	// Database configuration
	Config configuration.ConfigDatabase

	// Pointer to the db connection
	conn *mongo.Client
}

// Method to create a connection pool to the database
func (db *Database) connect() (*mongo.Client, error) {
	// build the connection url as a string
	// dsn := fmt.Sprintf(
	// 	"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	// 	db.Config.Host,
	// 	db.Config.Username,
	// 	db.Config.Password,
	// 	db.Config.Dbname,
	// 	db.Config.Port,
	// )

	// build the credential for the database authentication
	credential := options.Credential{
                Username: db.Config.Username,
                Password: db.Config.Password,
        }

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://"+db.Config.Host+":"+db.Config.Port).SetAuth(credential))
	errors.Raise(err)

	//set the command timeout
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    err = client.Connect(ctx)
    
    errors.Raise(err)
	
    // TO-DO: cleanup of the connection may be required in the longer run
    // defer client.Disconnect(ctx)
    
    // test of the database is reachable or not
    err = client.Ping(ctx, readpref.Primary())
    errors.Raise(err)

	db.conn = client

	// returns the connection pool as a pointer
	return db.conn, err
}

// Get the connection to the database or create a new one.
func (db *Database) Connection() *mongo.Client {
	// check if there is any connection at all
	if db.conn != nil {
		return db.conn
	} else {
		conn, err := db.connect()
		errors.Raise(err)
		return conn
	}
}
