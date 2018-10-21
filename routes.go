package main

import (
	"io/ioutil"
	"net/http"

	"github.com/posttul/lot-service/storage"
	"github.com/posttul/service"
	"github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
	yaml "gopkg.in/yaml.v2"
)

type routes = service.Routes

// Config is the service configuration
type Config struct {
	Service struct {
		Address string       `default:":3000"`
		Log     logrus.Level `default:"3"`
	}
	DB struct {
		User     string
		Password string
		DBName   string
		Host     string
	}
}

// Build gets all the config
func (c *Config) build(file string) error {
	config, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(config, c)
}

// LotService default service definition
type LotService struct {
	storage storage.Service
	routes  service.Routes
}

// GetRoutes returns the routes of the service
func (s LotService) GetRoutes() routes {
	s.routes = routes{"lot": {
		Path:    "/lot",
		Method:  http.MethodGet,
		Handler: s.GetLot(),
	}}
	return s.routes
}

// GetLot returns all lots on the storage
func (s *LotService) GetLot() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Retrive lots
		lots, err := s.storage.GetLots()
		if err != nil {
			service.Respose{Err: err}.Error(w)
			return
		}
		service.Respose{Data: lots}.OK(w)
	}
}
