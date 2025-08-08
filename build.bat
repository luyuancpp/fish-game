goctl rpc protoc ./apps/user/user.proto --zrpc_out=./apps/user --go_out=./apps/user --go-grpc_out=./apps/user -m
goctl rpc protoc ./apps/room/room.proto --zrpc_out=./apps/room --go_out=./apps/room --go-grpc_out=./apps/room -m

pause