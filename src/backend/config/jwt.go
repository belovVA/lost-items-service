package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type jwtConfig struct {
	jwt string `env:"JWT_TOKEN" env-required:"true"`
}

func (j *jwtConfig) GetSecret() string {
	return j.jwt
}

func JWTConfigLoad() (*jwtConfig, error) {
	var jwtCfg jwtConfig

	if err := cleanenv.ReadEnv(&jwtCfg); err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return &jwtCfg, nil

}
