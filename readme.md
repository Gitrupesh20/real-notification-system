## Overview

*A minimal real time notification system in go that delivers messages to connected clients over websockets using message broker rabbit mq,  
it uses worker base consumer to consumer message from queue and process them and broadcast message to specific user or to all connect users*



###  Key Highlights

- **Architecture:** Event-driven, Pub/Sub-style WebSocket fan-out
- **Transport:**
    - HTTP → For publishing messages
    - WebSocket → For delivering messages to clients
    - RabbitMQ → For message queuing & broker
- **Concurrency:**
    - Uses `sync.Map` to manage the WebSocket connection pool
    - Per-user buffered channel writer goroutines for async delivery
- **Performance:**
    - Handles **500+ messages/sec** without loss
    - Sub-**40ms latency** end-to-end


----

## Flow/ High Level Architecture

![img.png](img.png)

-----

## Running it Locally

### Prerequisites

* Go 1.20+
* RabbitMQ running locally and reachable at `mq_addr` in config (`amqp://guest:guest@localhost:5672/`)

---

### Steps

1. **Clone the repo**

   ```
   git clone https://github.com/yourusername/realtime-notification.git
   cd realtime-notification
   ```

2. **Ensure RabbitMQ is running**

   ```
   docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management
   ```

3. **Adjust configuration (if needed)**
   Update `cmd/config/config.example.json` to modify:

    * HTTP / WebSocket port
    * Queue name
    * Origins
    * SSL policy

4. **Run the server**

   ```
   go run ./cmd/main.go
   ```

5. **Connect a WebSocket client**

   ```
   ws://localhost:<port>/ws?userId=<your_user_id>
   ```

6. **Send a push notification**

   ```
   curl -X POST http://localhost:<port>/push \
     -H "Content-Type: application/json" \
     -d '{"user_id":"123", "message":"Hello 👋"}'
   ```

---

## Future Enhancements

* Add Prometheus metrics and Grafana dashboards
* JWT-based authentication for secure user sessions
* Message persistence via Redis or PostgreSQL
* Horizontal scaling with multiple RabbitMQ consumers

---



