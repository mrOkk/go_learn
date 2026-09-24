module LoadBalancer

go 1.27

replace ShortenerContract => ../ShortenerContract

require ShortenerContract v0.0.0

require (
	github.com/caarlos0/env/v11 v11.4.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/grpc v1.84.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
