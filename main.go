package main

import (
	"os"

	"github.com/posttul/lot-service/storage/postgres"
	"github.com/posttul/service"
	"github.com/sirupsen/logrus"
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
