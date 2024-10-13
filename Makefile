gen:
	protoc --go_out=. --go-grpc_out=. internal/grpc/*.proto

gen-user-service:
	protoc --go_out=. --go-grpc_out=. internal/grpc/user_service.proto

gen-trip-service:
	protoc --go_out=. --go-grpc_out=. internal/grpc/trip_service.proto

gen-payment-service:
	protoc --go_out=. --go-grpc_out=. internal/grpc/payment_service.proto

clean:
	rm internal/grpc/pb/*.go

# run-server:
# 	go run server/server.go server/models.go
# run-client:
# 	go run client/client.go