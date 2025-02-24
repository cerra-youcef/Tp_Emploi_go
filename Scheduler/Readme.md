# Scheduler for Fetching and Publishing Events

This app implements a scheduled task that periodically fetches resources from config api and correspending events from uca calendar api , then publishes those events to a NATS streaming service. It utilizes URLs stored in an `.env` file to fetch the resources and events.

## Features

- **Periodic Execution**: The scheduler runs every fixed number of seconds (adjustable) to fetch and publish events.
- **Event Publishing**: After fetching events, the scheduler publishes the events to a NATS streaming service.

## Setup
### Installation

Install the required Go dependencies:
```
go mod tidy
```

Run :
```
go run cmd/main.go
```

Subscribes to messages (events) that are published :
```
nats subscribe "EVENTS.>"
```
