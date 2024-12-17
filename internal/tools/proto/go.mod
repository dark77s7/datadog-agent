module github.com/DataDog/datadog-agent/internal/tools/proto

go 1.22.0

require (
	github.com/favadi/protoc-go-inject-tag v1.4.0
	github.com/golang/mock v1.5.0
	github.com/golang/protobuf v1.5.4
	github.com/grpc-ecosystem/grpc-gateway v1.12.2
	github.com/planetscale/vtprotobuf v0.4.0
	github.com/tinylib/msgp v1.1.6
	google.golang.org/grpc v1.24.0
)

require (
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/tools v0.26.0 // indirect
	google.golang.org/genproto v0.0.0-20190927181202-20e1ac93f88c // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace google.golang.org/protobuf v1.33.0 => google.golang.org/protobuf v1.34.0
