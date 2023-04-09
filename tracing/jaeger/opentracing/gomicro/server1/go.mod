module server1

go 1.16

require (
	github.com/go-micro/cli v1.1.4
	github.com/go-micro/plugins/v4/wrapper/trace/opentracing v1.2.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	go-micro.dev/v4 v4.9.0
	google.golang.org/protobuf v1.30.0
)

replace server1 => ./
