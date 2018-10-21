package main

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/posttul/lot-service/storage"
	"github.com/posttul/service"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type routes = service.Routes
type response struct {
	Status string
	Data   interface{}
	Error  error `yaml:"err,omitempty" json:"error,omitempty"`
}

func (r *response) SetStatus(status string) {
	r.Status = status
}

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

// YAML is a service.Writer to yaml
func YAML(w io.Writer, r service.Response) {
	b, err := yaml.Marshal(r)
	if err != nil {
		service.Error(w, &response{Error: err}, service.JSON)
		return
	}
	io.WriteString(w, string(b))
}

// GetLot returns all lots on the storage
func (s *LotService) GetLot() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Retrive lots
		lots, err := s.storage.GetLots()
		if err != nil {
			service.Error(w, &response{Error: err}, service.JSON)
			return
		}
		var outParser service.Writer
		switch r.URL.Query().Get("type") {
		case "yaml":
			outParser = YAML
		default:
			outParser = service.JSON
		}
		service.OK(w, &response{Data: lots}, outParser)
	}
}
