FROM golang:1.22.2-alpine AS builder

COPY . /github.com/ukrainskykirill/chat-server/source
WORKDIR /github.com/ukrainskykirill/chat-server/source

ARG DB_USER
ARG DB_PASSWORD
ARG DB_HOST
ARG DB_PORT
ARG DB_DATABASE_NAME
ARG GRPC_PORT

RUN go mod download
RUN go build -o ./bin/chat_server cmd/chat-server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/ukrainskykirill/chat-server/source/bin/chat_server .

ENV DB_USER=${DB_USER}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_DATABASE_NAME=${DB_DATABASE_NAME}
ENV GRPC_PORT=${GRPC_PORT}

CMD ["./chat_server"]