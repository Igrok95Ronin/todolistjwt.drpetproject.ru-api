package routes

import (
	"encoding/json"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/httperror"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type Check struct {
	Check bool `json:"check"`
}

// Отметить выполненную запись
func (h *handler) markCompletedEntry(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID := ps.ByName("id")

	var check Check

	if err := json.NewDecoder(r.Body).Decode(&check); err != nil {
		// Если произошла ошибка декодирования, возвращаем клиенту ошибку с кодом 400
		httperror.WriteJSONError(w, "Ошибка декодирования в JSON", err, http.StatusBadRequest)
		// Логируем ошибку
		h.logger.Errorf("Ошибка декодирования в JSON: %s", err)
		return
	}

	id, err := strconv.Atoi(ID)
	if err != nil {
		h.logger.Errorf("Некорректный ID: %s", err)
		return
	}
	if id <= 0 {
		h.logger.Errorf("ID должен быть больше 0: %d", id)
		return
	}

	if err = h.db.Model(&models.AllNotes{}).Where("id = ?", id).Update("completed", check.Check).Error; err != nil {
		h.logger.Errorf("Ошибка обновления записи с ID %d: %s", id, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
