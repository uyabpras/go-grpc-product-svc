proto:
	protoc -Ipkg/proto  --go_out=. --go-grpc_out=require_unimplemented_servers=false:. pkg/proto/*.proto


server:
	go run cmd/main.go