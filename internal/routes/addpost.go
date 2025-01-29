package routes

import (
	"encoding/json"
	"fmt"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/httperror"
	"github.com/julienschmidt/httprouter"
	"html"
	"net/http"
	"strings"
)

type Note struct {
	Note string `json:"note"`
}

// Добавить пост
func (h *handler) addPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.logger.Error("Не удалось получить user_id из контекста")
		httperror.WriteJSONError(w, "Не удалось получить user_id из контекста", nil, http.StatusInternalServerError)
		return
	}

	var note Note

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		// Если произошла ошибка декодирования, возвращаем клиенту ошибку с кодом 400
		httperror.WriteJSONError(w, "Ошибка декодирования в JSON", err, http.StatusBadRequest)
		// Логируем ошибку
		h.logger.Errorf("Ошибка декодирования в JSON: %s", err)
		return
	}

	note.Note = html.EscapeString(strings.TrimSpace(note.Note))
	if note.Note == "" {
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := addPostDB(h, w, note.Note, userID); err != nil {
		h.logger.Errorf("Ошибка при добавления записи в БД: %s", err)
	}
}

func addPostDB(h *handler, w http.ResponseWriter, note string, userID int64) error {

	allNotes := models.AllNotes{
		Note:   note,
		UserID: userID,
	}

	if err := h.db.Create(&allNotes).Error; err != nil {
		httperror.WriteJSONError(w, "Ошибка при добавления записи в БД", fmt.Errorf(""), http.StatusInternalServerError)
		return err
	}

	return nil
}
