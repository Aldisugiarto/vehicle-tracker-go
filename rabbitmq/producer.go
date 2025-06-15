package rabbitmq

import (
	"encoding/json"
	"log"
	"vehicle-tracker/models"

	"github.com/rabbitmq/amqp091-go"
)

var channel *amqp091.Channel
var conn *amqp091.Connection

func InitRabbitMQ() {
	var err error
	conn, err = amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ after retries: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	if err := ch.ExchangeDeclare("fleet.events", "direct", true, false, false, false, nil); err != nil {
		log.Fatal("Failed to declare exchange:", err)
	}
	if _, err := ch.QueueDeclare("geofence_alerts", true, false, false, false, nil); err != nil {
		log.Fatal("Failed to declare queue:", err)
	}
	if err := ch.QueueBind("geofence_alerts", "geofence_alerts", "fleet.events", false, nil); err != nil {
		log.Fatal("Failed to bind queue:", err)
	}
	channel = ch
}

func PublishGeofenceEvent(loc models.Location) {
	event := map[string]interface{}{
		"vehicle_id": loc.VehicleID,
		"event":      "geofence_entry",
		"location": map[string]float64{
			"latitude":  loc.Latitude,
			"longitude": loc.Longitude,
		},
		"timestamp": loc.Timestamp,
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Println("Failed to marshal event:", err)
		return
	}

	err = channel.Publish(
		"fleet.events",
		"geofence_alerts",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println("Failed to publish event:", err)
	}
}
