package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	//1) try to connect to rabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	//2) start listening for messages (interacting with the Q of rabbitMQ)
	log.Println("Listening for and Consuming RabbitMQ messages ...")

	//3) _create_ a consumer of q messages
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	//4) watch the queue. consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}

}

func connect() (*amqp.Connection, error) {

	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	//don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {

			fmt.Println("RabbitMQ not yet ready ...")
			counts++

		} else {

			log.Println("Connected to RabbitMQ!")
			connection = c
			break

		}

		if counts > 6 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off geometrically ...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil

}
