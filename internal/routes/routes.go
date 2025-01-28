package routes

import (
	"fmt"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/config"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/handlers"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var _ handlers.Handler = &handler{}

type handler struct {
	cfg    *config.Config
	logger *logging.Logger
}

func NewHandler(cfg *config.Config, logger *logging.Logger) handlers.Handler {
	return &handler{
		cfg:    cfg,
		logger: logger,
	}
}

func (h *handler) Router(router *httprouter.Router) {
	router.GET("/", h.Home)
}

func (h *handler) Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "TEST")
}
