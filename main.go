package main

import (
	"vehicle-tracker/controllers"
	"vehicle-tracker/database"
	"vehicle-tracker/mqtt"
	"vehicle-tracker/rabbitmq"

	"github.com/gin-gonic/gin"
)

func main() {
	csvFilePath := "mqtt/track.csv"
	vehicleID := "B1234XYZ"
	brokerHost := "mqtt"
	brokerPort := 1884

	database.InitDB()
	go rabbitmq.InitRabbitMQ()
	go mqtt.StartSubscriber()
	go mqtt.PublishCSV(csvFilePath, vehicleID, brokerHost, brokerPort)

	router := gin.Default()
	router.GET("/vehicles/:id/location", controllers.GetLastLocation)
	router.GET("/vehicles/:id/history", controllers.GetLocationHistory)

	router.Run(":8080")
	// Start the server on port 8080
}
