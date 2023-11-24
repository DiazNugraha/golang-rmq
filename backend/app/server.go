package app

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func Run() {

	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/interface")
	failOnerror(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnerror(err, "failed to open a channel")

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	failOnerror(err, "failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')

		err = ch.PublishWithContext(
			ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(text),
			},
		)

		failOnerror(err, "failed to publish a message")
		log.Printf("[x] Sent %s\n", text)
	}

}

func failOnerror(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}
