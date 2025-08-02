package configs

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Db    DbConfig   // postgres config
	Redis RdConfig   // redis config
	Auth  AuthConfig // auth jwt config
}

type DbConfig struct {
	Dsn string
}

type RdConfig struct {
	Addr string // redis addr from docker-compose
	Pass string // password
	NDB  string // BD number
}

type AuthConfig struct {
	SecretAcc       string // secret for access token
	SecretRef       string // secret for refresh token
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

const (
	timeExpAccessToken  = time.Minute * 15
	timeExpRefreshToken = time.Hour * 24
)

func LoadConfig() *Config {
	err := godotenv.Load("c:\\Son_Alex\\Go_projects\\e-commerce_proj\\VVV\\V2_mobile\\wisp\\backend\\simple_gin_server\\.env")
	if err != nil {
		fmt.Println("Error loading .env file, using default config", err.Error())
	}
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			SecretAcc:       os.Getenv("JWT_ACC_SECRET"),
			SecretRef:       os.Getenv("JWT_REF_SECRET"),
			AccessTokenExp:  timeExpAccessToken,
			RefreshTokenExp: timeExpRefreshToken,
		},
		Redis: RdConfig{
			Addr: os.Getenv("REDDIS_ADDR"),
			Pass: os.Getenv("REDDIS_PASS"),
			NDB:  os.Getenv("REDDIS_BD"),
		},
	}
}
