package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/posttul/lot-service/storage"
	"github.com/posttul/service"
	"github.com/sirupsen/logrus"
)

// ClientService service
type ClientService struct {
}

// GetRoutes gets client routes
func (cs *ClientService) GetRoutes() service.Routes {
	return service.Routes{"home": {
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			resp, err := http.Get("http://localhost:3000/")
			if err != nil {
				service.Error(w, &service.R{Error: err}, service.JSON)
				return
			}
			defer resp.Body.Close()
			lots := struct{ Data []storage.Lot }{}
			err = json.NewDecoder(resp.Body).Decode(&lots)
			if err != nil {
				service.Error(w, &service.R{Error: err}, service.JSON)
				return
			}
			service.OK(w, &service.R{Data: lots.Data}, service.JSON)
		},
	}}
}

func main() {
	service.SetLog(logrus.New())
	service.Start(":2000", &ClientService{})
}
