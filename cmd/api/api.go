package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justKody/ecom-go/service/cart"
	"github.com/justKody/ecom-go/service/order"
	"github.com/justKody/ecom-go/service/product"
	"github.com/justKody/ecom-go/service/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subRouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subRouter)

	orderStore := order.NewStore(s.db)

	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subRouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
