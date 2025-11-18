package rest

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"wallet/config"
	userHandler "wallet/rest/handlers/user"
	walletHandler "wallet/rest/handlers/walletB"
	middleware "wallet/rest/middlewares"
)

type Server struct {
	cnf           *config.Config
	userHandler   *userHandler.Handler
	walletHandler *walletHandler.Handler
}

func NewServer(
	cnf *config.Config,
	walletHdl *walletHandler.Handler,
	userHdl *userHandler.Handler,
) *Server {
	return &Server{
		cnf:           cnf,
		walletHandler: walletHdl,
		userHandler:   userHdl,
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

	server.walletHandler.RegisterRoutes(mux, manager)
	server.userHandler.RegisterRoutes(mux, manager)

	addr := ":" + strconv.Itoa(server.cnf.HttpPort)
	fmt.Println("Server running on", addr)
	err := http.ListenAndServe(addr, wrappedMux)
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
}
