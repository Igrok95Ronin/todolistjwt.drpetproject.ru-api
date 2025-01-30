package routes

import (
	"encoding/json"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/httperror"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Информация о пользователе
func (h *handler) me(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.logger.Error("Не удалось получить user_id из контекста")
		httperror.WriteJSONError(w, "Не удалось получить user_id из контекста", nil, http.StatusInternalServerError)
		return
	}

	var result struct {
		UserName string
	}

	if err := h.db.Table("users").Select("user_name").Where("id = ?", userID).Scan(&result).Error; err != nil {
		h.logger.Error("Ошибка при запросе к БД", err)
		httperror.WriteJSONError(w, "Ошибка при запросе к БД", nil, http.StatusInternalServerError)
		return
	}

	// Отправляем JSON-ответ с user_name
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
