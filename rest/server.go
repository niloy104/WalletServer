package rest

import (
	"wallet/config"

	"fmt"
	"net/http"
	"os"
	"strconv"
	"wallet/rest/handlers/user"
	middleware "wallet/rest/middlewares"
)

type Server struct {
	cnf *config.Config
	//productHandler *product.Handler
	userHandler *user.Handler
}

func NewServer(
	cnf *config.Config,
	//productHandler *product.Handler,
	userHandler *user.Handler,

) *Server {
	return &Server{
		cnf: cnf,
		//productHandler: productHandler,
		userHandler: userHandler,
	}
}

func (server *Server) Start() {
	manager := middleware.NewManager()
	manager.Use(
		middleware.Preflight,
		middleware.Cors,
		middleware.Logger,
	)
	mux := http.NewServeMux()
	wrappedMux := manager.WrapMux(mux)

	//server.productHandler.RegisterRoutes(mux, manager)
	server.userHandler.RegisterRoutes(mux, manager)

	addr := ":" + strconv.Itoa(server.cnf.HttpPort)
	fmt.Println("Server running on", addr)
	err := http.ListenAndServe(addr, wrappedMux)
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
}
