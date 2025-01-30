package routes

import (
	"encoding/json"
	"fmt"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/httperror"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

// Регистрация (создание нового пользователя)
// POST /register
func (h *handler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var users models.Users

	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		httperror.WriteJSONError(w, "Ошибка декодирования в json", err, http.StatusBadRequest)
		h.logger.Errorf("Ошибка декодирования в json: %s", err)
		return
	}

	username := strings.TrimSpace(users.UserName)
	email := strings.TrimSpace(users.Email)
	password := strings.TrimSpace(users.PasswordHash)

	// Проверка, что поля заполнены
	if username == "" || email == "" || password == "" {
		httperror.WriteJSONError(w, "Все поля (username, email, password) обязательны", nil, http.StatusConflict)
		h.logger.Errorf("Все поля (username, email, password) обязательны: %s", nil)
		return
	}

	//Запрет на выполнение скриптов
	username = template.HTMLEscapeString(username)
	email = template.HTMLEscapeString(email)
	password = template.HTMLEscapeString(password)

	// Проверка валидности email
	if err := ValidateEmail(email); err != nil {
		httperror.WriteJSONError(w, err.Error(), nil, http.StatusBadRequest)
		h.logger.Errorf("Неверный формат email: %s", email)
		return
	}

	// Проверяем, не существует ли уже такого пользователя
	//var existingUser User
	// Если запрос вернёт nil ошибку, значит пользователь найден, и он уже есть
	if err := h.db.Where("user_name = ? OR email = ?", username, email).First(&users).Error; err == nil {
		httperror.WriteJSONError(w, "Пользователь с таким username или email уже существует", err, http.StatusConflict)
		h.logger.Errorf("Пользователь с таким username или email уже существует: %s", err)
		return
	}

	// Хешируем пароль
	hashedPassword, err := HashPassword(password)
	if err != nil {
		httperror.WriteJSONError(w, "Ошибка при хешировании пароля", err, http.StatusInternalServerError)
		h.logger.Errorf("Ошибка при хешировании пароля: %s", err)
		return
	}

	// Создаём объект нового пользователя
	newUser := models.Users{
		UserName:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	// Сохраняем пользователя в базе
	if err = h.db.Create(&newUser).Error; err != nil {
		httperror.WriteJSONError(w, "Ошибка при сохранении пользователя", err, http.StatusInternalServerError)
		h.logger.Errorf("Ошибка при сохранении пользователя: %s", err)
		return
	}
	// Отправляем ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Пользователь успешно зарегистрирован"))
}

//---------------------------------------------------------------------------------------
//                                 УТИЛИТНЫЕ ФУНКЦИИ
//---------------------------------------------------------------------------------------

// HashPassword - хеширует пароль с помощью bcrypt (с cost = bcrypt.DefaultCost).
func HashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword вернёт хеш пароля.
	// bcrypt.DefaultCost по умолчанию равен 10 (можно увеличить, чтобы усложнить подбор).
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ValidateEmail(email string) error {
	// Проверка длины email
	if len(email) < 4 {
		return fmt.Errorf("email должен быть не менее 4 символов")
	}

	// Проверка наличия символа "@" в email
	if !strings.Contains(email, "@") {
		return fmt.Errorf("email должен содержать символ '@'")
	}

	// Проверка позиции символа "@" (не должен быть первым или последним символом)
	if strings.HasPrefix(email, "@") || strings.HasSuffix(email, "@") {
		return fmt.Errorf("email не может начинаться или заканчиваться на '@'")
	}

	// Дополнительно: базовая проверка формата email с помощью регулярного выражения
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return fmt.Errorf("ошибка проверки email: %v", err)
	}
	if !matched {
		return fmt.Errorf("email не соответствует формату")
	}

	return nil
}
