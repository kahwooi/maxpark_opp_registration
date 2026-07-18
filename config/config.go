package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type Config struct {
	Port                              string
	NATSUrl                           string
	OSSEndpoint                       string
	OSSAccessKeyID                    string
	OSSAccessKeySecret                string
	OSSBucketName                     string
	TurnstileSecretKey                string
	RegisterIndividualSubject         string
	RegisterIndividualFinalizeSubject string
	RegisterEmployerSubject           string
	RegisterEmployerFinalizeSubject   string
	RegisterEmployerIDSubject         string
	SiteCode                          string
	JWTSecretKey                      string
	MailerFrom                        string
	MailerHost                        string
	MailerPort                        int
	MaiilerUsername                   string
	MailerPassword                    string
}

func Load() *Config {
	env := os.Getenv("APP_ENV")
	envFile := ".env"
	if env != "" {
		envFile = fmt.Sprintf(".env.%s", env)
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("No .env file found, using system environment variables")
	}

	return &Config{
		Port:                              getEnv("PORT", "8081"),
		NATSUrl:                           getEnv("NATS_URL", ""),
		OSSEndpoint:                       getEnv("OSS_ENDPOINT", ""),
		OSSAccessKeyID:                    getEnv("OSS_ACCESS_KEY_ID", ""),
		OSSAccessKeySecret:                getEnv("OSS_ACCESS_KEY_SECRET", ""),
		OSSBucketName:                     getEnv("OSS_BUCKET_NAME", ""),
		TurnstileSecretKey:                getEnv("TURNSTILE_SECRET_KEY", ""),
		RegisterIndividualSubject:         getEnv("REGISTER_INDIVIDUAL_SUBJECT", ""),
		RegisterIndividualFinalizeSubject: getEnv("REGISTER_INDIVIDUAL_FINALIZE_SUBJECT", ""),
		RegisterEmployerSubject:           getEnv("REGISTER_EMPLOYER_SUBJECT", ""),
		RegisterEmployerFinalizeSubject:   getEnv("REGISTER_EMPLOYER_FINALIZE_SUBJECT", ""),
		RegisterEmployerIDSubject:         getEnv("REGISTER_EMPLOYER_ID_SUBJECT", ""),
		SiteCode:                          getEnv("SITE_CODE", ""),
		JWTSecretKey:                      getEnv("JWT_SECRET_KEY", ""),
		MailerFrom:                        getEnv("MAILER_FROM", ""),
		MailerHost:                        getEnv("MAILER_HOST", ""),
		MailerPort:                        getEnvAsInt("MAILER_PORT", 0),
		MaiilerUsername:                   getEnv("MAILER_USER", ""),
		MailerPassword:                    getEnv("MAILER_PASSWORD", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Error converting environment variable %s to int: %v, using default value %d", key, err, defaultValue)
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}

func InitializeMailer() *gomail.Dialer {
	mailerHost := os.Getenv("MAILER_HOST")
	mailerPortStr := os.Getenv("MAILER_PORT")
	mailerPort, err := strconv.Atoi(mailerPortStr)
	if err != nil {
		log.Fatalf("Error converting MAILER_PORT to int: %v", err)
	}
	mailerUsername := os.Getenv("MAILER_USER")
	mailerPassword := os.Getenv("MAILER_PASSWORD")

	return gomail.NewDialer(mailerHost, mailerPort, mailerUsername, mailerPassword)
}
