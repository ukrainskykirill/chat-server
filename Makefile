include .env

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(DB_PORT) dbname=$(DB_DATABASE_NAME) user=$(DB_USER) password=$(DB_PASSWORD) sslmode=disable"
POSTGRES_DSN="postgresql://postgres:postgres@localhost:5432/chat_server?sslmode=disable"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
	GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports@latest

lint:
	GOBIN=$(LOCAL_BIN) golangci-lint run ./... --config .golangci.pipeline.yaml

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-chat-api

generate-chat-api:
	mkdir -p pkg/chat_v1
	protoc --proto_path api/chat_v1 \
	--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/chat_v1/chat.proto

build:
	GOOS=linux GOARCH=amd64 go build -o chat_server_linux cmd/main.go

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t <REGESTRY>/chat-server:v0.0.1 .
	docker login -u <USERNAME> -p <PASSWORD> <REGESTRY>
	docker push <REGESTRY>/chat-server:v0.0.1 .

copy-to-server:
	scp chat_server_linux root@89.104.117.151:

install-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} -dbstring ${POSTGRES_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/ukrainskykirill/chat-server/internal/service/...,github.com/ukrainskykirill/chat-server/internal/api/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/ukrainskykirill/chat-server/internal/service/...,github.com/ukrainskykirill/chat-server/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

grpc-load-test:
	ghz \
		--proto api/chat_v1/chat.proto \
		--call chat_v1.ChatV1.Create \
		--data '{"userIDs": ["23","34"]}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051

grpc-error-load-test:
	ghz \
		--proto api/chat_v1/chat.proto \
		--call chat_v1.ChatV1.Create \
		--data '{"userIDs": ["1","2"]}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051

migrations-up:
	$(LOCAL_BIN)/goose -dir ./migrations postgres "postgresql://postgres:postgres@localhost:5432/chat_server?sslmode=disable" up