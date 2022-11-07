package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	mongoHost           = "MONGO_HOST"
	mongoDB             = "MONGO_DATABASE"
	mongoCredAuthMech   = "MONGO_CRED_AUTH_MECH"
	mongoCredAuthSource = "MONGO_CRED_AUTH_SOURCE"
	mongoCredUser       = "MONGO_CRED_USER"
	mongoCredPass       = "MONGO_CRED__PASSWORD"
	jwtSecret           = "JWT_SECRET"
	appPort             = "APP_PORT"
)

type MongoCred struct {
	AuthMechanism string `json:"authmechanism"`
	AuthSource    string `json:"authsource"`
	Username      string `json:"username"`
	Password      string `json:"password"`
}

type config struct {
	MongoHost string    `json:"mongohost"`
	MongoCred MongoCred `json:"mongocred"`
	MongoDB   string    `json:"mongodb"`
	JWTSecret string    `json:"jwtsecret"`
	AppPort   string    `json:"port"`
}

var filePath = "./config/config.json"

func Init() error {

	configFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("could not open or read config.json")
		return err
	}

	config := new(config)

	if err = json.Unmarshal(configFile, config); err != nil {
		log.Println("incorrect config file")
		return err
	}
	if err = os.Setenv(mongoHost, config.MongoHost); err != nil {
		log.Printf("can't set environment variable %s", mongoHost)
		return err
	}
	if err = os.Setenv(mongoCredAuthMech, config.MongoCred.AuthMechanism); err != nil {
		log.Printf("can't set environment variable %s", mongoCredAuthMech)
		return err
	}
	if err = os.Setenv(mongoCredAuthSource, config.MongoCred.AuthSource); err != nil {
		log.Printf("can't set environment variable %s", mongoCredAuthSource)
		return err
	}
	if err = os.Setenv(mongoCredUser, config.MongoCred.Username); err != nil {
		log.Printf("can't set environment variable %s", mongoCredUser)
		return err
	}

	if err = os.Setenv(mongoCredPass, config.MongoCred.Password); err != nil {
		log.Printf("can't set environment variable %s", mongoCredPass)
		return err
	}

	if err = os.Setenv(mongoDB, config.MongoDB); err != nil {
		log.Printf("can't set environment variable %s", mongoDB)
		return err
	}
	fmt.Println(mongoDB, config.MongoDB)

	if err = os.Setenv(jwtSecret, config.JWTSecret); err != nil {
		log.Printf("can't set environment variable %s", jwtSecret)
		return err
	}

	if err = os.Setenv(appPort, config.AppPort); err != nil {
		log.Printf("can't set environment variable %s", appPort)
		return err
	}

	return nil
}
