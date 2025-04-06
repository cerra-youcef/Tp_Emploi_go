# **UCA Study Planning Notification System**  
**Automated email alerts for changes in study planning at UCA University**  

## **📌 Overview**  
This system monitors changes in study planning (timetables, groups, grades) at UCA University and notifies subscribed users via email. It uses **NATS Streaming** for event-driven communication between microservices.  

---

## **🚀 Components**  

### **1. Config API**  
- **Purpose**: Stores resources (study groups, grades) and alert configurations (emails).  
- **Tech**: Go, REST API, SQLite.  

### **2. Timetable API**  
- **Purpose**:  
  - Stores current timetables.  
  - Compares events from the **Scheduler** and publishes alerts via NATS if changes detected (created/modified/deleted).  
- **Tech**: Go, NATS Streaming, SQLite.  
- **Consumes**: `events` (NATS).  
- **Publishes**: `alerts` (NATS).  

### **3. Scheduler**  
- **Purpose**: Fetches events from UCA’s server and forwards them to the **Timetable Service** via NATS.  
- **Tech**: Go, NATS Streaming.  
- **Publishes**: `events` (NATS).  

### **4. Alerter**  
- **Purpose**:  
  - Listens to `alerts` (NATS).  
  - Queries **Config API** for emails.  
  - Sends notifications (emails).  
- **Tech**: Go, NATS, GCC (Garde Champêtre Consulting) for emails.  

### **5. NATS Server**  
- **Role**: Message broker for inter-service communication.  
- **Channels**:  
  - `EVENTS` (Scheduler → Timetable).  
  - `ALERTS` (Timetable → Alerter).  

---

## **🔧 API Documentation**  
### **Config API (Swagger)**  
Run the server and visit:  
`http://localhost:8080/swagger/index.html`  

### **Timetable API (Swagger)**  
Run the server and visit:  
`http://localhost:8081/swagger/index.html`  

