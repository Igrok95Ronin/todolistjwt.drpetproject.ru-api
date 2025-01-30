package routes

import (
	"encoding/json"
	"fmt"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/httperror"
	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

// RefreshHandler - обработчик обновления токенов.
// Использует refresh-токен для выдачи нового access-токена и нового refresh-токена.
func (h *handler) Refresh(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Извлекаем refresh_token из куки
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		httperror.WriteJSONError(w, "Необходим refresh_token (cookie отсутствует)", err, http.StatusUnauthorized)
		h.logger.Errorf("Необходим refresh_token (cookie отсутствует): %s", err)
		return
	}
	refreshToken := refreshCookie.Value

	// 2. Валидируем refresh-токен
	claims, err := ValidateRefreshToken(h, refreshToken)
	if err != nil {
		httperror.WriteJSONError(w, "Невалидный или просроченный refresh-токен", err, http.StatusUnauthorized)
		h.logger.Errorf("Невалидный или просроченный refresh-токен: %s", err)
		return
	}

	// 3. Проверим, что пользователь существует в базе
	var users models.Users
	if err := h.db.First(&users, claims.UserID).Error; err != nil {
		httperror.WriteJSONError(w, "Пользователь не найден", err, http.StatusUnauthorized)
		h.logger.Errorf("Пользователь не найден: %s", err)
		return
	}

	// 4. Проверим, что refresh-токен совпадает с тем, что хранится в базе
	if users.RefreshToken != refreshToken {
		httperror.WriteJSONError(w, "Refresh-токен не соответствует сохранённому в базе", err, http.StatusUnauthorized)
		h.logger.Errorf("Refresh-токен не соответствует сохранённому в базе: %s", err)
		return
	}

	// 5. Генерируем новые токены
	newAccessToken, err := GenerateAccessToken(h, users.ID)
	if err != nil {
		httperror.WriteJSONError(w, "Ошибка при генерации нового access-токена", err, http.StatusInternalServerError)
		h.logger.Errorf("Ошибка при генерации нового access-токена: %s", err)
		return
	}
	newRefreshToken, err := GenerateRefreshToken(h, users.ID)
	if err != nil {
		httperror.WriteJSONError(w, "Ошибка при генерации нового refresh-токена", err, http.StatusInternalServerError)
		h.logger.Errorf("Ошибка при генерации нового refresh-токена: %s", err)
		return
	}

	// 6. Сохраняем новый refresh-токен в базе
	users.RefreshToken = newRefreshToken
	if err = h.db.Save(&users).Error; err != nil {
		httperror.WriteJSONError(w, "Ошибка при сохранении нового refresh-токена", err, http.StatusInternalServerError)
		h.logger.Errorf("Ошибка при сохранении нового refresh-токена: %s", err)
		return
	}

	// 7. Обновляем куки (access и refresh)
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	// 8. Успешный ответ с информацией о токенах
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}
	json.NewEncoder(w).Encode(response)
}

// ValidateRefreshToken - парсит и валидирует refresh-токен. Возвращает claims, если успешно.
func ValidateRefreshToken(h *handler, refreshToken string) (*models.MyClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Возвращаем секрет, которым подписан refresh-токен
		return []byte(h.cfg.Token.Refresh), nil
	}

	// Парсим токен
	parsedToken, err := jwt.ParseWithClaims(refreshToken, &models.MyClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	// Проверяем, что claims верного типа и токен валиден
	claims, ok := parsedToken.Claims.(*models.MyClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("Невалидный токен")
	}

	return claims, nil
}
