module github.com/dundduun/msg/sso

go 1.25.7

require (
	buf.build/go/protovalidate v1.1.1
	github.com/dundduun/msg/protos v0.0.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.3
	github.com/ilyakaznacheev/cleanenv v1.5.0
	go.uber.org/zap v1.27.1
	google.golang.org/grpc v1.78.0
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260209202127-80ab13bee0bf.1 // indirect
	cel.dev/expr v0.25.1 // indirect
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/google/cel-go v0.27.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20260209203927-2842357ff358 // indirect
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260209200024-4cfbd4190f57 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260209200024-4cfbd4190f57 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/dundduun/msg/protos => ../protos
