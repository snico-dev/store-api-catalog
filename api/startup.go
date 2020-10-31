package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NicolasDeveloper/store-catalog-api/api/controllers"
	"github.com/NicolasDeveloper/store-catalog-api/infra"
	"github.com/golobby/container"
	"github.com/gorilla/mux"
)

type config struct {
	port string
}

type startup struct {
	config config
	router *mux.Router
}

//StartUp inilialize api
func StartUp(port string) error {
	s := &startup{
		config: config{
			port,
		},
	}

	infra.NewContainer()

	s.registerDb()
	s.registerRoutes()
	s.run()

	return nil
}

func (s *startup) registerRoutes() {
	s.router = mux.NewRouter()
	subrouter := s.router.PathPrefix("/catalog-api/v1/").Subrouter()

	subrouter.HandleFunc("/health-check", controllers.HealthCheck).Methods(http.MethodGet)
	subrouter.HandleFunc("/products", controllers.CreateProduct).Methods(http.MethodPost)
	subrouter.HandleFunc("/products", controllers.UpdateProduct).Methods(http.MethodPut)
	subrouter.HandleFunc("/products/active", controllers.ActiveProduct).Methods(http.MethodPut)
	subrouter.HandleFunc("/products/inactive", controllers.InactiveProduct).Methods(http.MethodPut)

	http.Handle("/", s.router)
}

func (s *startup) registerDb() {
	var dbctx *infra.DbContext
	container.Make(&dbctx)
	dbctx.Connect()
}

func (s *startup) run() {
	fmt.Printf("API Catalog Running on %v", s.config.port)
	log.Fatal(http.ListenAndServe(":"+s.config.port, s.router))
}
