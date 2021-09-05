package database

import (
	"time"
	"context"
	"testing"

	"github.com/riotpot/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/riotpot/internal/database"
	"github.com/riotpot/internal/configuration"
)

// Runs a simple test to check the health of MongoDB database and the configuration.
// Keep in mind this function requires the role `superuser` and the `superuser` db
// to previously exists!
// Example:
//  $ # create the user as a superuser and password as password under admin db
//  $ db.createUser( { user: "superuser", pwd: "password", roles: ["root" ] } )

func TestDatabaseConnection(t *testing.T) {

	var (
		// Load an identity for the database
		id = configuration.ConfigIdentity{
			Name: "test_db",
		}

		// Load a configuration for the database
		config = configuration.ConfigDatabase{
			Username: "superuser",
			Password: "password",
			Host:     "127.0.0.1",
			Identity: id,
			Port:     "27017",
		}

		// Load a database object
		db = database.Database{
			Config: config,
		}

	)

	// Connect to the db
	conn := db.Connection()

	if conn == nil{
		t.Error("Error connecting database")
	}
	defer conn.Disconnect(context.TODO())
}

func TestDatabaseInsert(t *testing.T) {

	var (
		// Load an identity for the database
		id = configuration.ConfigIdentity{
			Name: "test_db",
		}

		db_name = "test_db"
		collection_name = "test_col"
		payload = "Test Run"
		// to store query output
		out []bson.M

		// Load a configuration for the database
		config = configuration.ConfigDatabase{
			Username: "superuser",
			Password: "password",
			Host:     "127.0.0.1",
			Identity: id,
			Port:     "27017",
		}

		// Load a database object
		db = database.Database{
			Config: config,
		}
	)

	// Connect to the db
	conn := db.Connection()
	// create a Test connection item to store
	test_model := models.TestConnection(payload)
	
	input_time := test_model.Timestamp

	if conn != nil{
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// create database
		db_element := conn.Database(db_name)
		// create collection
		collec_element := db_element.Collection(collection_name)
		// // store item in MongoDB
		_, err := collec_element.InsertOne(ctx, test_model)
		
		if err != nil {
            t.Error(err)
        }

        // retrieve item from MongoDB
		res, err := collec_element.Find(ctx, bson.D{})

		if err = res.All(ctx, &out); err != nil {
 		   t.Error(err)
		}

		len_slice := len(out)
		got_time := out[len_slice-1]["timestamp"]
		got_payload	:= out[len_slice-1]["payload"]

		// check if the correct item is picked
		if got_time != input_time {
			t.Error(err)
		}

		// check if the payload match
		if got_payload != payload {
			t.Error(err)	
		}

		// // cleanup
 		if err = db_element.Drop(ctx); err != nil {
			t.Error(err)
		}

		defer cancel()
		defer conn.Disconnect(context.TODO())
	} else {
		t.Error("Database not accessible")
	}
}
