package config

type jwtConfig struct {
	Secret         string `yaml:"secret" json:"secret"`
	RefreshExpTime uint   `yaml:"refresh_exp" json:"refresh_exp"`
	AccessExpTime  uint   `yaml:"access_exp" json:"access_exp"`
}

func defaultJwtConfig() jwtConfig {
	return jwtConfig{
		Secret:         "mysecret",
		RefreshExpTime: 30,
		AccessExpTime:  15,
	}
}

func (p *jwtConfig) loadFromEnv() {
	loadEnvStr("JWT_SECRET", &p.Secret)
	loadEnvUint("JWT_REFRESH_TOKEN_EXP_TIME", &p.RefreshExpTime)
	loadEnvUint("JWT_ACCESS_TOKEN_EXP_TIME", &p.AccessExpTime)
}
