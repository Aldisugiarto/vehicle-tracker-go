package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"regexp"
	"time"
	"vehicle-tracker/database"
	"vehicle-tracker/models"
	"vehicle-tracker/rabbitmq"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func distanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// Validasi payload
func validateLocation(loc models.Location) error {
	if loc.VehicleID == "" || len(loc.VehicleID) < 3 {
		return fmt.Errorf("invalid vehicle_id")
	}
	if loc.Latitude < -90 || loc.Latitude > 90 {
		return fmt.Errorf("invalid latitude")
	}
	if loc.Longitude < -180 || loc.Longitude > 180 {
		return fmt.Errorf("invalid longitude")
	}
	if loc.Timestamp <= 0 {
		return fmt.Errorf("invalid timestamp")
	}
	return nil
}

// Handler untuk pesan masuk
var MessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	log.Printf("Received message on topic: %s", topic)

	matched, _ := regexp.MatchString(`^/fleet/vehicle/[^/]+/location$`, topic)
	if !matched {
		log.Println("Topic tidak valid, abaikan")
		return
	}

	var loc models.Location
	if err := json.Unmarshal(msg.Payload(), &loc); err != nil {
		log.Println("Payload bukan JSON yang valid:", err)
		return
	}

	if err := validateLocation(loc); err != nil {
		log.Println("Validasi gagal:", err)
		return
	}

	parsedTime := time.Unix(loc.Timestamp, 0)
	log.Printf("Valid Location: %+v at %s\n", loc, parsedTime.Format(time.RFC3339))

	if distanceMeters(loc.Latitude, loc.Longitude, -6.1754, 106.8272) <= 50 {
		rabbitmq.PublishGeofenceEvent(loc)
	}

	err := insertLocationToDB(loc)
	if err != nil {
		log.Println("Failed to insert location into database:", err)
		return
	}

	log.Println("Location successfully inserted into database")
}

func insertLocationToDB(loc models.Location) error {
	query := `
        INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp)
        VALUES ($1, $2, $3, $4)
    `
	_, err := database.DB.Exec(query, loc.VehicleID, loc.Latitude, loc.Longitude, time.Unix(loc.Timestamp, 0))
	return err
}

func StartSubscriber() {
	broker := "tcp://mqtt:1884"
	clientID := "fleet-location-subscriber"

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetDefaultPublishHandler(MessageHandler)

	client := mqtt.NewClient(opts)

	for i := 0; i < 5; i++ { // Retry up to 5 times
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Printf("MQTT connect failed (attempt %d): %v", i+1, token.Error())
			time.Sleep(2 * time.Second)
		} else {
			log.Println("MQTT connected successfully")
			break
		}
	}

	topic := "/fleet/vehicle/+/location"
	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatal("MQTT subscribe failed:", token.Error())
	}

	log.Println("MQTT subscriber started on topic:", topic)

	select {}
}
