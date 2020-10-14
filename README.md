# Message Queue Service
### Golang, Echo, AMQP

#### Setup
* Config database and amqp config
* Config service mode and destination:
    ```yaml
    ts_service:
      mode:
      destination: order
    ```
    * Mode will be 0: run publisher and consumer, 1: run publisher, 2: run consumer. If you want to run both, keep it blank.
    * Destination is name of destination service, it will decide url of publisher.
        Example if destination is `order`, url will be `http://localhost:9090/ems/api/v1/order/messages`
* Install require packages: `go mod vendor`

#### Startup
* Run: `go run main.go`
* Publish messages:
```
curl --location --request POST 'localhost:9090/ems/v1/order/messages' \
--header 'Content-Type: application/json' \
--data-raw '{
    "routing_key": "routing.key",
    "payload": {
        "name": "This is message"
    },
    "origin_code": "CODE",
    "origin_model": "MODEL"
}'
```
| Fields       | Type          | Required | Not Null | Description                       |
|:-------------|:-------------:|:--------:|:--------:|:----------------------------------|
| routing_key  | string        | YES      | YES      | Routing key                       |
| payload      | json          | YES      | YES      | Message content (json)            |
| origin_model | string        | NO       | NO       | Object model                      |
| origin_code  | string        | NO       | NO       | Object code                       |
| headers      | list string   | NO       | NO       | Headers will be sent with message |

#### Documents:
See documents at: http://localhost:9090/ems/v1/{ts_service.destination}/swagger/index.html  
For example: http://localhost:9090/ems/v1/order/swagger/index.html
![](https://i.imgur.com/IUxywZy.png)

#### Diagram
![alt text](https://i.imgur.com/KwUNR1V.png)


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
- [AMQP Golang](https://godoc.org/github.com/streadway/amqp)

#### Contributing
If you want to contribute to this boilerplate, clone the repository and just start making pull requests.
