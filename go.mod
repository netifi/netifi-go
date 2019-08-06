module github.com/netifi/netifi-go

go 1.12

require (
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/hoisie/mustache v0.0.0-20120318181656-6dfe7cd5e765
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/rsocket/rsocket-go v0.3.2
	github.com/rsocket/rsocket-rpc-go v0.0.1
	github.com/stretchr/testify v1.3.0
)

replace github.com/rsocket/rsocket-go => ../rsocket-go

replace github.com/rsocket/rsocket-rpc-go => ../rsocket-rpc-go

replace github.com/panjf2000/ants => ../ants
