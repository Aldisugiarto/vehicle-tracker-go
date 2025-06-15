package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func checkRabbitMQConnection(host string, port string) error {
	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("RabbitMQ is not reachable at %s: %v", address, err)
	}
	defer conn.Close()
	return nil
}

func main() {
	rabbitHost := "rabbitmq"
	rabbitPort := "5672"

	// Check RabbitMQ connection before dialing
	if err := checkRabbitMQConnection(rabbitHost, rabbitPort); err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	conn, err := amqp091.Dial(fmt.Sprintf("amqp://guest:guest@%s:%s/", rabbitHost, rabbitPort))
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"geofence_alerts",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to register consumer:", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg map[string]interface{}
			if err := json.Unmarshal(d.Body, &msg); err == nil {
				fmt.Printf("📍 Geofence Alert Received: %v\n", msg)
			} else {
				log.Println("Failed to parse message:", err)
			}
		}
	}()

	fmt.Println(" [*] Waiting for geofence alerts. To exit press CTRL+C")
	<-forever
}
