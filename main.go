package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"api-telemetria-robo/controller"
	"api-telemetria-robo/logs"
	"api-telemetria-robo/repository"
	"api-telemetria-robo/service"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var pkgName logs.PackageName = "MAIN"

func main() {
	client, err := connectMongoDBClient("mongodb+srv://<username>:<password>@cluster0.jiue3nk.mongodb.net/?appName=Cluster0")
	if err != nil {
		log.Fatal(err.Error())
	}
	matchRepository := repository.NewMatchMongo(client, "test", "matches")

	matchService := service.NewMatchService(matchRepository)
	roundService := service.NewRoundService(matchRepository)

	matchController := controller.NewMatchController(matchService, roundService)

	mux := mux.NewRouter()
	matchController.LoadRoutes(mux)

	serverAddress := fmt.Sprintf("%s:%s", "0.0.0.0", "8080")
	server := http.Server{
		Addr:    serverAddress,
		Handler: mux,
	}

	logs.Infof(pkgName, "Starting server at %s", serverAddress)
	server.ListenAndServe()
}

func connectMongoDBClient(uri string) (*mongo.Client, error) {
	var (
		client *mongo.Client
		err    error
	)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err = mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
