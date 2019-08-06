module github.com/netifi/netifi-go

go 1.12

require (
	github.com/google/uuid v1.1.1
	github.com/jjeffcaii/reactor-go v0.0.12
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/rsocket/rsocket-go v0.0.0-00010101000000-000000000000
	github.com/rsocket/rsocket-rpc-go v0.0.1
	github.com/stretchr/testify v1.3.0
)

replace github.com/rsocket/rsocket-go => ../rsocket-go

replace github.com/rsocket/rsocket-rpc-go => ../rsocket-rpc-go

replace github.com/panjf2000/ants => ../ants
