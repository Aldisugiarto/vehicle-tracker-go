package controllers

import (
	"database/sql"
	"net/http"
	"vehicle-tracker/database"
	"vehicle-tracker/models"

	"github.com/gin-gonic/gin"
)

func GetLastLocation(c *gin.Context) {
	vehicleID := c.Param("id")
	var loc models.Location
	var timestamp sql.NullTime

	row := database.DB.QueryRow(
		"SELECT vehicle_id, latitude, longitude, timestamp FROM vehicle_locations WHERE vehicle_id = $1 ORDER BY timestamp DESC LIMIT 1",
		vehicleID,
	)

	err := row.Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &timestamp)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if timestamp.Valid {
		loc.Timestamp = timestamp.Time.Unix()
	}

	c.JSON(http.StatusOK, loc)
}

func GetLocationHistory(c *gin.Context) {
	vehicleID := c.Param("id")
	start := c.Query("start")
	end := c.Query("end")
	var timestamp sql.NullTime

	query := "SELECT vehicle_id, latitude, longitude, timestamp FROM vehicle_locations WHERE vehicle_id = $1"
	args := []interface{}{vehicleID}

	if start != "" && end != "" {
		query += " AND timestamp BETWEEN TO_TIMESTAMP($2) AND TO_TIMESTAMP($3) ORDER BY timestamp ASC"
		args = append(args, start, end)
	} else {
		query += " ORDER BY timestamp ASC"
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var history []models.Location
	for rows.Next() {
		var loc models.Location
		if err := rows.Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &timestamp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if timestamp.Valid {
			loc.Timestamp = timestamp.Time.Unix()
		}
		history = append(history, loc)
	}
	if history == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no history found"})
		return
	}

	c.JSON(http.StatusOK, history)
}
