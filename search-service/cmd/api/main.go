package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"search-service/internal/data"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort    = ":80"
	mongoURL   = "mongodb://mongo:27017"
	ctxTimeOut = 15 * time.Second
)

var client *mongo.Client

func main() {
	var err error

	// Service flags
	mongoUsername := flag.String("mongoUsername", "", "MongoDB user")
	mongoPassword := flag.String("mongoPassword", "", "MongoDB user password")

	flag.Parse()

	if *mongoUsername == "" || *mongoPassword == "" {
		log.Fatal("MongoDB username or password cannot be empty. Exiting...")
	}

	// connecting to mongo
	client, err = connectToMongo(*mongoUsername, *mongoPassword)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	closeMongoDBConn := func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal("Error disconnecting from mongo:", err)
		}
	}

	defer closeMongoDBConn()

	// Giving the mongodb connection to the data package
	data.NewConn(client)

	// starting web server
	srv := &http.Server{
		Addr:    webPort,
		Handler: routes(),
	}

	log.Println("Starting SearchService on port", webPort)

	err = srv.ListenAndServe()

	log.Fatal(err)
}

// connectToMongo establishes a mongodb connvetion
// and returns a *mongo.Client or an error
func connectToMongo(username, password string) (*mongo.Client, error) {
	// creating connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: username,
		Password: password,
	})

	// connect
	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to mongo:", err)
		return nil, err
	}

	return conn, nil
}
