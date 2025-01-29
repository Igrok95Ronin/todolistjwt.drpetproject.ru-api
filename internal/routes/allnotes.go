package routes

import (
	"encoding/json"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/httperror"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *handler) allNotes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := r.Context().Value("user_id").(int64)

	var allNotes []models.AllNotes

	if err := h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&allNotes).Error; err != nil {
		h.logger.Error(err)
		httperror.WriteJSONError(w, "Ошибка при получения данных из бд", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(allNotes); err != nil {
		h.logger.Error(err)
		httperror.WriteJSONError(w, "Ошибка при отправке данных клиенту", err, http.StatusInternalServerError)
	}
}
