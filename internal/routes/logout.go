package routes

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

// LogoutHandler - обработчик логаута
func (h *handler) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Устанавливаем куки с прошедшей датой
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0), // просрочен
		HttpOnly: true,
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Вы успешно вышли из системы"))
}
