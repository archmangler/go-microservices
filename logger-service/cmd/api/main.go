package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017" //later change localhost -> mongo ... for some reason
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {

	log.Println("logger service entering main() ...") //debugging

	//connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	//create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	//Register the RPC server to the App
	err = rpc.Register(new(RPCServer))

	//start the RPC server on the logger service
	go app.rpcListen()

	go app.gRPCListen()

	//start the web server
	log.Println("(main) Starting service on port: ", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Println("(main) Listening and Serving on port: ", webPort)

	err = srv.ListenAndServe()

	if err != nil {
		log.Println("(main) failed to start service on port: ", webPort)
		log.Panic()
	}

}

func (app *Config) rpcListen() error {

	log.Println("Starting RPC server on port ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {

	log.Println("logger service entering connectToMongo ...") //debugging

	//create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//actually connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("(connectToMongo) Error connecting:", err)
		return nil, err
	}

	log.Println("(connectToMongo) Connected to mongodb ...")

	return c, nil
}
