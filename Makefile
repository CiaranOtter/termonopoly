all:
	protoc --go_out=. --go-grpc_out=. proto/game.proto
# --go_out=comm  --go-grpc_out=comm
# --go_opt=paths=source_relative
# --go-grpc_opt=paths=source_relative
clean:
	rm -rf ./termonopoly/comm/* 