package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/posttul/lot-service/storage"
	"github.com/posttul/service"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type routes = service.Routes
type response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  interface{} `yaml:"err,omitempty" json:"error,omitempty"`
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

// NotFound handle the default error not not finding anything
func (s *LotService) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	service.JSON(w, &response{Error: "not found", Status: "error"})
}

// InitRouter use modify the behavior of the router in the service.
func (s *LotService) InitRouter(r *httprouter.Router) *httprouter.Router {
	r.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		service.JSON(w, &response{Error: "not allow", Status: "error"})
	})
	r.NotFound = http.HandlerFunc(s.NotFound)
	return r
}

// GetRoutes returns the routes of the service.
func (s *LotService) GetRoutes() routes {
	s.routes = routes{"lot": {
		Path:    "/lot/:id",
		Method:  http.MethodGet,
		Handler: s.GetLot(),
	}, "alllot": {
		Path:    "/lot",
		Method:  http.MethodGet,
		Handler: s.GetAllLots(),
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

// GetAllLots returns all the lots from the storage.
func (s *LotService) GetAllLots() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

// GetLot returns all lots on the storage
func (s *LotService) GetLot() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Retrive lots
		if p.ByName("id") == "" {
			log.Debug("Request lot with missing id")
			s.NotFound(w, r)
			return
		}
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			log.Errorf("Error on string to int conversion %s", err.Error())
			s.NotFound(w, r)
			return
		}
		lot, err := s.storage.GetLotByID(int64(id))
		if err != nil {
			log.Debugf("Error on GetLotByID %s", err.Error())
			s.NotFound(w, r)
			return
		}
		service.OK(w, &response{
			Data: lot,
		}, service.JSON)

	}
}
