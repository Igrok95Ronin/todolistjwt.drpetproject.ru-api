package routes

import (
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// Удалить запись
func (h *handler) deleteEntry(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID := ps.ByName("id")

	id, err := strconv.Atoi(ID)
	if err != nil {
		h.logger.Errorf("Некорректный ID: %s", err)
		return
	}
	if id <= 0 {
		h.logger.Errorf("ID должен быть больше 0: %d", id)
		return
	}

	if err = h.db.Where("id = ?", id).Unscoped().Delete(&models.AllNotes{}).Error; err != nil {
		h.logger.Errorf("Ошибка удаления записи с ID %d: %s", id, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
