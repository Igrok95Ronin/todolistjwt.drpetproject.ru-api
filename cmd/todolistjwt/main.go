package main

import (
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/config"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/routes"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func main() {
	db := routes.InitDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database handle: %v", err)
	}
	defer sqlDB.Close()

	router := httprouter.New()

	logger := logging.GetLogger()

	cfg := config.GetConfig()

	handler := routes.NewHandler(cfg, logger, db)
	handler.Router(router)

	corsHandler := routes.CorsSettings().Handler(router)

	start(corsHandler, cfg, logger)
}

func start(router http.Handler, cfg *config.Config, logger *logging.Logger) {

	const WRI = 15 * time.Second

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		WriteTimeout: WRI,
		ReadTimeout:  WRI,
		IdleTimeout:  WRI,
	}

	logger.Infof("Server started %v", cfg.Port)
	logger.Fatal(server.ListenAndServe())
}
