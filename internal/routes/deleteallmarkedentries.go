package routes

import (
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/httperror"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Удалить все отмеченные записи
func (h *handler) DeleteAllMarkedEntries(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.logger.Error("Не удалось получить user_id из контекста")
		httperror.WriteJSONError(w, "Не удалось получить user_id из контекста", nil, http.StatusInternalServerError)
		return
	}

	if err := h.db.Unscoped().Where("completed = ? AND user_id = ?", true, userID).Delete(&models.AllNotes{}).Error; err != nil {
		h.logger.Errorf("Ошибка при удалении всех отмеченных записей: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
