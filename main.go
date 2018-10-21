package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"bitbucket.org/housing/lot/service"
	"bitbucket.org/housing/lot/storage/postgres"
)

var log = logrus.New()

func main() {

	c := &Config{}
	err := c.build("./config.yaml")
	if err != nil {
		panic(err)
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	service.SetLog(log)

	sto, err := postgres.New(c.DB.User, c.DB.Password, c.DB.DBName, c.DB.Host)
	if err != nil {
		panic(err)
	}

	service.Start(
		c.Service.Address,
		LotService{
			storage: sto,
		})
}
