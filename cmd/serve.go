package cmd

import (
	"fmt"
	"os"
	"wallet/config"
	"wallet/infra/db"
	"wallet/repo"
	"wallet/rest"
	"wallet/user"
	"wallet/wallet"

	usrHandler "wallet/rest/handlers/user"
	walletHandler "wallet/rest/handlers/walletB"
	middleware "wallet/rest/middlewares"
)

func Serve() {
	// Load config
	cnf := config.GetConfig()

	// Database connection
	dbCon, err := db.NewConnection(cnf.DB)
	if err != nil {
		fmt.Println("DB connection error:", err)
		os.Exit(1)
	}

	// Run migrations
	if err := db.MigrateDB(dbCon, "./migrations"); err != nil {
		fmt.Println("Migration error:", err)
		os.Exit(1)
	}

	// Repositories
	walletRepo := repo.NewWalletRepo(dbCon)
	userRepo := repo.NewUserRepo(dbCon)

	// DB Executor for transactions
	dbExec := repo.NewDBExecutor(dbCon)

	// Services
	usrSvc := user.NewService(userRepo)
	walletSvc := wallet.NewService(walletRepo, dbExec)

	// Middlewares
	mws := middleware.NewMiddlewares(cnf)

	// Handlers
	userHdl := usrHandler.NewHandler(cnf, usrSvc)
	walletHdl := walletHandler.NewHandler(walletSvc, mws) // service first, middlewares second

	// Start server
	server := rest.NewServer(cnf, walletHdl, userHdl)
	server.Start()
}
