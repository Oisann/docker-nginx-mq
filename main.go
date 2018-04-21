package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func reload() {
	command := exec.Command("kill", "-HUP", "$NGINX_PID")
	command.Dir = "/"
	out, err := command.Output()
	failOnError(err, "Reload failed!")
	fmt.Println("Produced: " + string(out))
}

func main() {
	conn, err := amqp.Dial(os.Getenv("AMQP_SERVER")) // example: amqp://rabbitmq:rabbitmq@rabbit1:5672/
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"command", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			if string(d.Body) == "reload" {
				go reload()
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
