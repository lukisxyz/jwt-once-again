package config

import "fmt"

type pgConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     uint   `yaml:"port" json:"port"`
	DBName   string `yaml:"db_name" json:"db_name"`
	SslMode  string `yaml:"ssl_mode" json:"ssl_mode"`
	Password string `yaml:"password" json:"password"`
	Username string `yaml:"user" json:"user"`
}

func (p *pgConfig) ConnStr() string {
	return fmt.Sprintf("host=%s port=%d database=%s sslmode=%s user=%s password=%s", p.Host, p.Port, p.DBName, p.SslMode, p.Username, p.Password)
}

func defaultPgConfig() pgConfig {
	return pgConfig{
		Host:     "localhost",
		Port:     5433,
		DBName:   "postgres",
		SslMode:  "disable",
		Password: "password",
		Username: "postgres",
	}
}

func (p *pgConfig) loadFromEnv() {
	loadEnvStr("DB_HOST", &p.Host)
	loadEnvUint("DB_PORT", &p.Port)
	loadEnvStr("DB_NAME", &p.DBName)
	loadEnvStr("DB_SSL", &p.SslMode)
	loadEnvStr("DB_PASSWORD", &p.Password)
	loadEnvStr("DB_USER", &p.Username)
}
