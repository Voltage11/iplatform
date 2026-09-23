package config

import (
	"time"
)

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func newJwtConfig() (*JWTConfig, error) {
	secret, err := getEnvReq("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	accessTTLInt, err := getEnvInt("JWT_ACCESS_TTL_SECONDS", 15*60)
	if err != nil {
		return nil, err
	}
	accessTTLDuration := time.Duration(accessTTLInt) * time.Second

	refreshTTLInt, err := getEnvInt("JWT_REFRESH_TTL_HOURS", 24*30)
	if err != nil {
		return nil, err
	}
	refreshTTLDuration := time.Duration(refreshTTLInt) * time.Hour

	return &JWTConfig{
		Secret:     secret,
		AccessTTL:  accessTTLDuration,
		RefreshTTL: refreshTTLDuration,
	}, nil
}
