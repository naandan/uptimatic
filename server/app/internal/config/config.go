package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppDebug    bool
	AppPort     string
	AppLogLevel string
	AppDomain   string
	AppScheme   string

	AuthJWTSecret              string
	AuthAccessTokenExpiration  time.Duration
	AuthRefreshTokenExpiration time.Duration

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	RedisHost string
	RedisPort int
	RedisPass string
	RedisDB   int

	EmailSmtpHost string
	EmailSmtpPort int
	EmailSmtpUser string
	EmailSmtpPass string
	EmailFrom     string

	GoogleClientId     string
	GoogleClientSecret string
}

func LoadConfig() (Config, error) {
	var cfg Config

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		viper.AutomaticEnv()
	}

	AuthAccessTokenExpiration, _ := time.ParseDuration(viper.GetString("AUTH_ACCESS_TOKEN_EXPIRATION"))
	AuthRefreshTokenExpiration, _ := time.ParseDuration(viper.GetString("AUTH_REFRESH_TOKEN_EXPIRATION"))

	cfg = Config{
		AppDebug:    viper.GetBool("APP_DEBUG"),
		AppPort:     viper.GetString("APP_PORT"),
		AppLogLevel: viper.GetString("APP_LOG_LEVEL"),
		AppDomain:   viper.GetString("APP_DOMAIN"),
		AppScheme:   viper.GetString("APP_SCHEME"),

		AuthJWTSecret:              viper.GetString("AUTH_JWT_SECRET"),
		AuthAccessTokenExpiration:  AuthAccessTokenExpiration,
		AuthRefreshTokenExpiration: AuthRefreshTokenExpiration,

		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetInt("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASS"),
		DBName:     viper.GetString("DB_NAME"),

		RedisHost: viper.GetString("REDIS_HOST"),
		RedisPort: viper.GetInt("REDIS_PORT"),
		RedisPass: viper.GetString("REDIS_PASS"),
		RedisDB:   viper.GetInt("REDIS_DB"),

		EmailSmtpHost: viper.GetString("EMAIL_SMTP_HOST"),
		EmailSmtpPort: viper.GetInt("EMAIL_SMTP_PORT"),
		EmailSmtpUser: viper.GetString("EMAIL_SMTP_USER"),
		EmailSmtpPass: viper.GetString("EMAIL_SMTP_PASS"),
		EmailFrom:     viper.GetString("EMAIL_FROM"),

		GoogleClientId:     viper.GetString("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
	}

	return cfg, nil
}

func (c *Config) DBDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func (c *Config) RedisURL() string {
	if c.RedisPass != "" {
		return fmt.Sprintf("redis://:%s@%s:%d/0", c.RedisPass, c.RedisHost, c.RedisPort)
	}
	return fmt.Sprintf("redis://%s:%d/0", c.RedisHost, c.RedisPort)
}
