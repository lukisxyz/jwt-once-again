package config

import "fmt"

type listenConfig struct {
	Host         string `yaml:"host" json:"host"`
	Port         uint   `yaml:"port" json:"port"`
	ReadTimeout  uint   `yaml:"read_to" json:"read_to"`
	WriteTimeout uint   `yaml:"write_to" json:"write_to"`
	IdleTimeout  uint   `yaml:"idle_to" json:"idle_to"`
}

func (l listenConfig) Addr() string {
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}

func defaultListenConfig() listenConfig {
	return listenConfig{
		Host:         "127.0.0.1",
		Port:         8080,
		ReadTimeout:  25,
		WriteTimeout: 25,
		IdleTimeout:  300,
	}
}

func (l *listenConfig) loadFromEnv() {
	loadEnvStr("LISTEN_HOST", &l.Host)
	loadEnvUint("LISTEN_PORT", &l.Port)
	loadEnvUint("LISTEN_READ_TIMEOUT", &l.ReadTimeout)
	loadEnvUint("LISTEN_WRITE_TIMEOUT", &l.WriteTimeout)
	loadEnvUint("LISTEN_IDLE_TIMEOUT", &l.IdleTimeout)
}
