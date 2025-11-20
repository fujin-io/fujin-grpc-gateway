package config

import "github.com/fujin-io/fujin/public/config"

type Config struct {
	Addr string           `yaml:"addr"`
	TLS  config.TLSConfig `yaml:"tls"`
	GRPC GRPCConfig       `yaml:"grpc"`
}

type GRPCConfig struct {
	Addr string           `yaml:"addr"`
	TLS  config.TLSConfig `yaml:"tls"`
	// TODO: Add more GRPC config here
}
