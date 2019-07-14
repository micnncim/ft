.PHONY: proto
proto:
	protoc --go_out=plugins=grpc:. proto/*.proto

.PHONY: server
server:
	go build -o bin/server github.com/micnncim/ft/server

.PHONY: client
client:
	go build -o bin/client github.com/micnncim/ft/client
