package cmd

import (
	"fmt"
	"os"
	"wallet/config"
	"wallet/infra/db"
	"wallet/repo"
	"wallet/rest"
	"wallet/user"

	usrHandler "wallet/rest/handlers/user"
)

func Serve() {
	cnf := config.GetConfig()

	//fmt.Println("%+v", cnf.DB)

	dbCon, err := db.NewConnection(cnf.DB)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = db.MigrateDB(dbCon, "./migrations")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// repos
	//productRepo := repo.NewProductRepo(dbCon)
	userRepo := repo.NewUserRepo(dbCon)

	// domains
	usrSvc := user.NewService(userRepo)
	//prdctSvc := product.NewService(productRepo)

	//middlewares := middleware.NewMiddlewares(cnf)

	// handlers
	//productHandler := prdctHandler.NewHandler(middlewares)
	userHandler := usrHandler.NewHandler(cnf, usrSvc)

	server := rest.NewServer(
		cnf,
		//productHandler,
		userHandler,
	)
	server.Start()
}
