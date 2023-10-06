.PHONY:run
run:
	go run cmd/main.go
.PHONY:gen
gen:
	protoc --go_out=gen --go-grpc_out=require_unimplemented_servers=false:gen image_service.proto 