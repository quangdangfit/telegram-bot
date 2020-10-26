# Telegram Bot Service
### Golang, Gin, Telegram Bot

#### Setup
* Config database and amqp config
* Install require packages: `go mod vendor`

#### Startup
* Run: `go run main.go`
* Publish messages:
```
curl --location --request POST 'http://localhost:9090/internal/messages' \
--header 'Authorization: Basic ZW1zOmVtc2ludGVybmFs' \
--header 'Content-Type: application/json' \
--data-raw '{
    "code": "SESSION_ABC1",
    "action": "add_session",
    "data": {
        "sessionCode": "SESSION_ABC1",
        "partnerTransportationCode": "T200903VTK7U5PR",
        "locationType": "GHN_HUB",
        "locationId": "2388",
        "stopIndex": 1,
        "type": "PICKUP"
    }
}'
```
#### Structure
```
├── app  
│   ├── api             # Handle request & response
│   ├── dbs             # Database Layer
│   ├── models          # Models
│   ├── queue           # AMQP Layer
│   ├── repositories    # Repositories Layer
│   ├── router          # Router api v1  
│   ├── schema          # Sechemas  
│   ├── services        # Business Logic Layer  
├── config              # Config's files 
├── docs                # Swagger API documents
├── pkg                 # Common packages
│   └── utils           # Utilities
```

#### 📙 Libraries
- [Gin](https://godoc.org/github.com/gin-gonic/gin)
- [Telegram API](https://core.telegram.org/bots/api)

#### Contributing
If you want to contribute to this boilerplate, clone the repository and just start making pull requests.
