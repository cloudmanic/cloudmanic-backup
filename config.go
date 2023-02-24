package main

import (
	"github.com/caarlos0/env"
	_ "github.com/jpfuentes2/go-env/autoload"
)

var (
	cfg Config
)

type Config struct {
	BackupName             string `env:"BACKUP_NAME"`
	AlertEmail             string `env:"ALERT_EMAIL"`
	ObjectRegion           string `env:"OBJECT_REGION"`
	ObjectBucket           string `env:"OBJECT_BUCKET"`
	ObjectAccessKeyId      string `env:"OBJECT_ACCESS_KEY_ID"`
	ObjectSecreteAccessKey string `env:"OBJECT_SECRET_ACCESS_KEY"`
	ObjectEndpoint         string `env:"OBJECT_ENDPOINT"`
	MysqlHost              string `env:"MYSQL_HOST" envDefault:"127.0.0.1"`
	MysqlPort              string `env:"MYSQL_PORT" envDefault:"3306"`
	MysqlDb                string `env:"MYSQL_DB"`
	MysqlUser              string `env:"MYSQL_USER"`
	MysqlPassword          string `env:"MYSQL_PASSWORD"`
	EncryptKey             string `env:"ENCRYPT_KEY"`
	BackupDbStoreDir       string `env:"BACKUP_DB_STORE_DIR"`
	HoursBetweenBackups    string `env:"HOURS_BETWEEN_BACKUPS"`
	PintSuccessUrl         string `env:"PING_SUCCESS_URL"`
	DbSizeCheckLow         int64  `env:"DB_SIZE_CHECK_LOW"`
	DbSizeCheckHigh        int64  `env:"DB_SIZE_CHECK_HIGH"`
	MailDriver             string `env:"MAIL_DRIVER"`
	MailHost               string `env:"MAIL_HOST"`
	MailPort               int    `env:"MAIL_PORT"`
	MailUsername           string `env:"MAIL_USERNAME"`
	MailPassword           string `env:"MAIL_PASSWORD"`
	MailEncryption         string `env:"MAIL_ENCRYPTION"`
	MailFromEmail          string `env:"MAIL_FROM_EMAIL"`
	MailGunDomain          string `env:"MAILGUN_DOMAIN"`
	MailGunApiKey          string `env:"MAILGUN_API_KEY"`
}

//
// Init....
//
func init() {
	cfg = Config{}

	// Parse the config
	err := env.Parse(&cfg)

	if err != nil {
		Log(err.Error())
	}
}

/* End File */
