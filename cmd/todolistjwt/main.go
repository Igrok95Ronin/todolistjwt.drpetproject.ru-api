package main

import (
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/config"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/routes"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func main() {
	router := httprouter.New()

	logger := logging.GetLogger()

	cfg := config.GetConfig()

	handler := routes.NewHandler(cfg, logger)
	handler.Router(router)

	start(router, cfg, logger)
}

func start(router *httprouter.Router, cfg *config.Config, logger *logging.Logger) {

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
