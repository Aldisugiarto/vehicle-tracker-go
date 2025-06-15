package mqtt

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"vehicle-tracker/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func PublishCSV(filePath string, vehicleID string, brokerHost string, brokerPort int) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", brokerHost, brokerPort))
	opts.SetClientID("mqtt_publisher")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	defer client.Disconnect(250)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}

	for _, record := range records {
		if len(record) < 2 {
			log.Println("Skipping invalid record:", record)
			continue
		}

		latitude := parseFloat(record[0])
		longitude := parseFloat(record[1])
		timestamp := time.Now().Unix()

		location := models.Location{
			VehicleID: vehicleID,
			Latitude:  latitude,
			Longitude: longitude,
			Timestamp: timestamp,
		}
		payload, err := json.Marshal(location)
		if err != nil {
			log.Printf("Failed to marshal location: %v", err)
			continue
		}

		topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicleID)
		token := client.Publish(topic, 0, false, payload)
		token.Wait()

		if token.Error() != nil {
			log.Printf("Failed to publish message: %v", token.Error())
		} else {
			log.Printf("Published: %s", string(payload))
		}

		time.Sleep(2 * time.Second)
	}
}

func parseFloat(value string) float64 {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Failed to parse float value: %v", err)
		return 0
	}
	return result
}
