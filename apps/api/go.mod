module github.com/seonNoh/seonology-journey/apps/api

go 1.24.0

toolchain go1.24.3

require (
	github.com/MicahParks/keyfunc/v3 v3.3.10
	github.com/coder/websocket v1.8.13
	github.com/go-chi/chi/v5 v5.2.1
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/seonNoh/seonology-journey/proto/gen/go v0.0.0
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/MicahParks/jwkset v0.8.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
)

replace github.com/seonNoh/seonology-journey/proto/gen/go => ../../proto/gen/go
