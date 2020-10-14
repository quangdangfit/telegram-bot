module transport/ems

go 1.14

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gin-gonic/gin v1.6.3
	github.com/go-openapi/spec v0.19.9 // indirect
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/google/uuid v1.1.1
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/mailru/easyjson v0.7.3 // indirect
	github.com/manucorporat/try v0.0.0-20170609134256-2a0c6b941d52
	github.com/spf13/viper v1.7.0
	github.com/streadway/amqp v1.0.0
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	go.uber.org/dig v1.10.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/tools v0.0.0-20200811215021-48a8ffc5b207 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/yaml.v2 v2.3.0 // indirect
	transport/lib v1.0.38
)

replace transport/lib v1.0.38 => gitlab.ghn.vn/logistics/ts/lib.git v1.0.38
