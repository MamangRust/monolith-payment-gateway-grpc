migrate:
	go run service/migrate/main.go up

migrate-down:
	go run service/migrate/main.go down


generate-proto:
	protoc --proto_path=pkg/proto --go_out=shared/pb --go_opt=paths=source_relative --go-grpc_out=shared/pb --go-grpc_opt=paths=source_relative pkg/proto/*.proto


generate-swagger:
	swag init -g service/apigateway/cmd/main.go -o service/apigateway/docs

seeder:
	go run service/seeder/main.go


api-gateway:
	go run service/apigateway/cmd/main.go


auth-service:
	go run service/auth/cmd/main.go

role-service:
	go run service/role/cmd/main.go

card-service:
	go run service/card/cmd/main.go

merchant-service:
	go run service/merchant/cmd/main.go

user-service:
	go run service/user/cmd/main.go

saldo-service:
	go run service/saldo/cmd/main.go

topup-service:
	go run service/topup/cmd/main.go

transaction-service:
	go run service/transaction/cmd/main.go


transfer-service:
	go run service/transfer/cmd/main.go


withdraw-service:
	go run service/withdraw/cmd/main.go

email-service:
	go run service/email/cmd/main.go

