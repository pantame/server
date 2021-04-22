package config

import (
	"encoding/json"
	"net/mail"
	"os"
)

type DBConfig struct {
	Host string `json:"host"`
	Port uint `json:"port"`
	Pass string `json:"pass"`
	User string `json:"user"`
	Db   string `json:"db"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	DB       int `json:"db"`
}

type PathsConfig struct {
	SaveDir string `json:"save_dir"`
}

type AccessTokensConfig struct {
	SendGridKey string `json:"send_grid_key"`
}

type EmailAddressesConfig struct {
	FromSupport mail.Address `json:"from_support"`
	FromNoReply mail.Address `json:"from_no_reply"`
	ToAdmin     mail.Address `json:"to_admin"`
}

type config struct {
	Database     DBConfig `json:"database"`
	Redis        RedisConfig `json:"redis"`
	Paths        PathsConfig `json:"paths"`
	AccessTokens AccessTokensConfig `json:"access_tokens"`
	Emails       EmailAddressesConfig `json:"emails"`
}

var conf config

func Read() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(file, &conf); err != nil {
		panic(err)
	}
}

func Database() DBConfig {
	return conf.Database
}

func Redis() RedisConfig {
	return conf.Redis
}

func Paths() PathsConfig {
	return conf.Paths
}

func AccessTokens() AccessTokensConfig {
	return conf.AccessTokens
}

func Emails() EmailAddressesConfig {
	return conf.Emails
}
