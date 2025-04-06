# Alerter
## **Email Alert Service**  
**Purpose**: Automatically sends email alerts by:  
1. Subscribing to the `ALERTS` channel on NATS.  
2. Fetching alert rules from the Config API.  
3. Triggering emails based on configured conditions.  
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
