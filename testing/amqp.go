package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"

	protobuf "code.google.com/p/goprotobuf/proto"
	"github.com/sterops/tinkerings/android-remote/proto"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"websocket", // name
		false,    // durable
		false,    // delete when usused
		false,    // exclusive
		false,    // noWait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	cmd := new(proto.Command)
	cmd.Type = proto.Command_COMMAND.Enum()
	cmd.Command = protobuf.String("Hello")
	data, err := protobuf.Marshal(cmd)
	failOnError(err, "Failed to marshal")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        data,
		})
	failOnError(err, "Failed to publish a message")

	os.Exit(0)
}
