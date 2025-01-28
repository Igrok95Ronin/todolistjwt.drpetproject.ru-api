package routes

import (
	"encoding/json"
	"fmt"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/config"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/handlers"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"net/http"
)

var _ handlers.Handler = &handler{}

type handler struct {
	cfg    *config.Config
	logger *logging.Logger
	db     *gorm.DB
}

func NewHandler(cfg *config.Config, logger *logging.Logger, db *gorm.DB) handlers.Handler {
	return &handler{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

func (h *handler) Router(router *httprouter.Router) {
	router.GET("/", h.Home)
}

func (h *handler) Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var allNotes []models.AllNotes

	if err := h.db.Find(&allNotes).Error; err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json") // Сначала заголовки
	w.WriteHeader(http.StatusOK)                       // Затем статус

	if err := json.NewEncoder(w).Encode(&allNotes); err != nil {
		fmt.Println(err)
	}
}
