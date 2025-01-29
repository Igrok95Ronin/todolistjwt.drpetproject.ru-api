package routes

import (
	"encoding/json"
	"fmt"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/httperror"
	"github.com/julienschmidt/httprouter"
	"html"
	"net/http"
	"strconv"
	"strings"
)

type ModifiedEntry struct {
	ModEntry string `json:"modEntry"`
}

// Редактировать запись
func (h *handler) editEntry(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var modifiedEntry ModifiedEntry

	if err := json.NewDecoder(r.Body).Decode(&modifiedEntry); err != nil {
		// Если произошла ошибка декодирования, возвращаем клиенту ошибку с кодом 400
		httperror.WriteJSONError(w, "Ошибка декодирования в JSON", err, http.StatusBadRequest)
		// Логируем ошибку
		h.logger.Errorf("Ошибка декодирования в JSON: %s", err)
		return
	}

	id, _ := strconv.Atoi(ps.ByName("id"))

	modifiedEntry.ModEntry = html.EscapeString(strings.TrimSpace(modifiedEntry.ModEntry))
	if modifiedEntry.ModEntry == "" {
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := editEntryDB(h, w, modifiedEntry.ModEntry, id); err != nil {
		h.logger.Errorf("Ошибка при обновлении записи по id: %v %s", id, err)
	}
}

func editEntryDB(h *handler, w http.ResponseWriter, editEntry string, id int) error {

	if err := h.db.Model(&models.AllNotes{}).Where("id = ?", id).Update("note", editEntry).Error; err != nil {
		httperror.WriteJSONError(w, "Ошибка при обновления записи в БД", fmt.Errorf(""), http.StatusInternalServerError)
		return err
	}

	return nil
}
