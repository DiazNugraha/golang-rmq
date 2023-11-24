package main

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/interface")
	failAndError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failAndError(err, "Failed to open a channel")
	defer ch.Close()
	for {
		msg, err := ch.Consume(
			"hello", // queue
			"",      // consumer
			true,    // auto-ack
			false,   // exclusive
			false,   // no-local
			false,   // no-wait
			nil,     // args
		)
		failAndError(err, "Failed to register a consumer")
		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		for d := range msg {
			log.Printf("Received a message: %s", d.Body)
		}
	}
}

func failAndError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
