## Run

Tidy / download modules :
```
go mod tidy
```
Build & run :
```
go run cmd/main.go
```

---
Or build : 
```
go build -o config cmd/main.go
```
Then run : 
```
./config
```

## Documentation

Documentation is visible in **api** directory.

or with Swagger UI after running the api: 
```
http://127.0.0.1:8080/swagger/index.html
```

