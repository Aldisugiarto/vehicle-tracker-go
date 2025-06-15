# Vehicle Tracking System

This project is a backend system to track vehicle locations in real-time, trigger geofence alerts, and expose a REST API to query location data.

## Features

- Receive real-time vehicle location data via MQTT
- Store location data in PostgreSQL
- Check if a vehicle enters a defined geofence (50m radius)
- Publish geofence alerts to RabbitMQ
- Worker service to consume and process geofence alerts
- Expose REST API to query last location and location history

---

## Technologies Used

- Golang
- Gin (web framework)
- PostgreSQL
- MQTT (Eclipse Mosquitto)
- RabbitMQ
- Docker & Docker Compose

---

## Project Structure

```bash
vehicle-tracker/
├── Dockerfile
├── Dockerfile.worker
├── docker-compose.yml
├── main.go
├── config/             # Environment config (optional)
├── db/                 # Database initialization
├── models/             # Data structures
├── handlers/           # API endpoint handlers
├── mqtt/               # MQTT subscriber and geofence logic
├── rabbitmq/           # RabbitMQ producer logic
├── geofence/worker.go  # Worker that listens for geofence alerts
```

## Running the project
1. Run docker compose in root directory (there is docker-compose file)
   ```
   docker-compose up --build
   ```
2. Check PostgreSQL Database via DBeaver (Optional)
   - Host: localhost
   - Port: 5433
   - Database: vehicle-tracker
   - Username: postgres
   - Password: postgres
3. REST-API 
   - Get Last location
   ```
   curl http://localhost:8080/vehicles/B1234XYZ/location
   ```
   - History (start - end)
   ```
   curl http://localhost:8080/vehicles/B1234XYZ/history?start=1714521600&end=1717804800
   ```
   - Other alternatif you can use Postman with collection in this roject
4. Check event alert in RabbitMQ UI
   - Open: http://localhost:15672
   - Username/Password: guest / guest

5. Geofence logic
   - Geofence is a circular zone with a 50-meter radius around:
     - Latitude: -6.1754
     - Longitude: 
