package util

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	AppName                string `mapstructure:"APP_NAME"`
	Environment            string `mapstructure:"ENVIRONMENT"`
	RedisAddress           string `mapstructure:"REDIS_ADDRESS"`
	WorkerCount            int    `mapstructure:"WORKER_COUNT"`
	MailSMTPServer         string `mapstructure:"MAIL_SMTP_SERVER"`
	MailSMTPServerPort     int    `mapstructure:"MAIL_SMTP_SERVER_PORT"`
	MailFromAddress        string `mapstructure:"MAIL_FROM_ADDRESS"`
	MailSMTPServerPassword string `mapstructure:"MAIL_SMTP_SERVER_PASSWORD"`
	MailSenderTitle        string `mapstructure:"MAIL_SENDER_TITLE"`
	MailDisableSend        bool   `mapstructure:"MAIL_DISABLE_SEND"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
