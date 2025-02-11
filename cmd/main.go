package main

import (
	"flag"
	"hangry/db"
	repo "hangry/db"
	"hangry/generated"
	"hangry/handler"
	"hangry/seeder"
	"hangry/usecase"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	dbDsn := os.Getenv("DATABASE_URL")
	db := db.NewDBConn(db.NewDBConnOptions{
		Dsn: dbDsn,
	})

	flag.Parse()
	args := flag.Args()

	// to handle seeding
	if len(args) >= 1 {
		switch args[0] {
		case "seed":
			seeder.Seed(db)
			os.Exit(0)
		}
	}

	// initiate repository
	cartRepo := repo.NewCartRepository(db)
	productRepo := repo.NewProductRepository(db)
	promoRepo := repo.NewPromoRepository(db)
	orderRepo := repo.NewOrderRepository(db)
	userRepo := repo.NewUserRepository(db)

	// this repo is for managing transaction
	transactionRepo := repo.NewTransactionRepository(db)

	// initiate usecase
	cartUsecase := usecase.NewCartUsecase(
		transactionRepo,
		cartRepo,
		productRepo,
	)
	promoUsecase := usecase.NewPromoUsecase(promoRepo, transactionRepo, cartRepo, productRepo)
	orderUsecase := usecase.NewOrderUsecase(transactionRepo, orderRepo, userRepo, cartRepo, promoRepo)

	var server generated.ServerInterface = handler.NewServer(
		cartUsecase,
		promoUsecase,
		orderUsecase,
	)

	generated.RegisterHandlers(e, server)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":1323"))
}
