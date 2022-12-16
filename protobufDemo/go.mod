module demo

go 1.19

require (
	github.com/golang/protobuf v1.5.2
	google.golang.org/protobuf v1.28.1

	"server" v0.0.0
)

replace (
	"server" => "../server"
)
