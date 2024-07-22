FROM golang:1.22.2-alpine AS builder

COPY . /github.com/ukrainskykirill/chat-server/source
WORKDIR /github.com/ukrainskykirill/chat-server/source


RUN go mod download
RUN go build -o ./bin/chat_server cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/ukrainskykirill/chat-server/source/bin/chat_server .

CMD ["./chat_server"]