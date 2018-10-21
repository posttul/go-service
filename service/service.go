package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

// Service main definition of the service
type Service interface {
	GetRoutes() Routes
}

type logHandler struct {
	handler http.Handler
}
type responseWriterLoger struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriterLoger) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (l *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	wl := &responseWriterLoger{w, http.StatusOK}
	l.handler.ServeHTTP(wl, r)
	log.WithFields(logrus.Fields{
		"Method": r.Method,
		"Path":   r.URL,
		"Status": wl.statusCode,
		"Time":   time.Since(t),
	}).Info("Request")
}

// Start is use to serve a service
func Start(addrs string, s Service) {
	router := httprouter.New()
	for _, r := range s.GetRoutes() {
		fmt.Printf("Register endpoint %s with the method %s and handler %T \n", r.Path, r.Method, r.Handler)
		router.Handle(r.Method, r.Path, r.Handler)
	}
	fmt.Printf("Starting server at port: %s and service %T ", addrs, s)
	err := http.ListenAndServe(addrs, &logHandler{router})
	if err != nil {
		panic(err)
	}
}

// Routes a slice of route
type Routes = map[string]Route

// Route on the service
type Route struct {
	Path    string
	Method  string
	Handler httprouter.Handle
}

// Respose simple envelop
type Respose struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Err    error       `json:"error,omitempty"`
}

func (r Respose) json(w io.Writer) {
	bts, err := json.Marshal(r)
	if err != nil && log != nil {
		log.Error(err)
	}
	io.WriteString(w, string(bts))
}

// OK write the ok response
func (r Respose) OK(w io.Writer) {
	r.Status = "ok"
	if hw, ok := w.(http.ResponseWriter); ok {
		hw.WriteHeader(http.StatusOK)
		r.json(hw)
		return
	}
	r.json(w)

}

// Error write the error response
func (r Respose) Error(w io.Writer) {
	r.Status = "error"
	if hw, ok := w.(http.ResponseWriter); ok {
		hw.WriteHeader(http.StatusInternalServerError)
		r.json(hw)
		return
	}
	r.json(w)
}

// Deny write the deny response
func (r Respose) Deny(w io.Writer) {
	r.Status = "deny"
	if hw, ok := w.(http.ResponseWriter); ok {
		hw.WriteHeader(http.StatusUnauthorized)
		r.json(hw)
		return
	}
	r.json(w)
}
