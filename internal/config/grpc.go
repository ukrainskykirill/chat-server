package config

import (
	"net"
	"os"
)

const (
	grpcHostEnv = "GRPC_HOST"
	grpcPortEnv = "GRPC_PORT"
)

type GRPCConfig interface {
	Address() string
}

type grpcConfig struct {
	host string
	port string
}

func NewGRPCConfig() (GRPCConfig, error) {
	host, ok := os.LookupEnv(grpcHostEnv)
	if !ok {
		return &grpcConfig{}, errVariableNotFound
	}

	port, ok := os.LookupEnv(grpcPortEnv)
	if !ok {
		return &grpcConfig{}, errVariableNotFound
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
