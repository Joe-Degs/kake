package kake

import "net"

type Config struct {
	Addr *net.TCPAddr
}

const DefaultAddr = "localhost:6996"

func DefaultConfig() *Config {
	cfg := &Config{}

	// get default tcp address, if its not available settle for a random port
	if err := cfg.ListenDefault(); err != nil {
		cfg.useRandoPort()
	}

	return cfg
}

func (c *Config) ListenDefault() error { return c.ListenAddress(DefaultAddr) }

func (c *Config) useRandoPort() {
	_ = c.ListenAddress("localhost:0")
}

func (c *Config) ListenAddress(addr string) (err error) {
	c.Addr, err = net.ResolveTCPAddr("tcp", addr)
	return
}
