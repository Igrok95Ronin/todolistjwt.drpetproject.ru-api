package routes

import (
	"fmt"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/config"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/handlers"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"net/http"
)

var _ handlers.Handler = &handler{}

type handler struct {
	cfg    *config.Config
	logger *logging.Logger
	db     *gorm.DB
}

func NewHandler(cfg *config.Config, logger *logging.Logger, db *gorm.DB) handlers.Handler {
	return &handler{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

func (h *handler) Router(router *httprouter.Router) {
	// Регистрация (создание нового пользователя)
	// POST /register
	router.POST("/register", h.Register)

	// Логин (получение access и refresh токенов)
	// POST /login
	router.POST("/login", h.Login)

	// Обновление (refresh) токенов
	//POST /refresh
	router.POST("/refresh", h.Refresh)

	// Защищённый маршрут, доступный только при наличии валидного access-токена
	// GET /protected
	//
	// Оборачиваем ProtectedHandler в AuthMiddleware.
	// Благодаря этому любой маршрут, который мы обернём AuthMiddleware,
	// станет защищённым, и проверка access-токена будет выполняться автоматически.
	router.GET("/protected", AuthMiddleware(h.Protected))
	// Выход из системы
	router.POST("/logout", h.Logout)

	router.GET("/", AuthMiddleware(h.allNotes))                                        // Получения всех записей
	router.POST("/addpost", AuthMiddleware(h.addPost))                                 // Добавить пост
	router.PUT("/editentry/:id", AuthMiddleware(h.editEntry))                          // Редактировать запись
	router.DELETE("/deleteentry/:id", AuthMiddleware(h.deleteEntry))                   // Удалить запись
	router.PUT("/markcompletedentry/:id", AuthMiddleware(h.markCompletedEntry))        // Отметить выполненную запись
	router.DELETE("/deleteallentries", AuthMiddleware(h.DeleteAllEntries))             // Удалить все записи
	router.DELETE("/deleteallmarkedentries", AuthMiddleware(h.DeleteAllMarkedEntries)) // Удалить все отмеченные записи
	router.GET("/me", AuthMiddleware(h.me))                                            // Информация о пользователе
}

// ProtectedHandler - обработчик примера защищённого маршрута.
// Доступ сюда возможен только через AuthMiddleware.
func (h *handler) Protected(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//// Достаём user_id из контекста (установили в AuthMiddleware).
	//userID, ok := r.Context().Value("user_id").(uint)
	//if !ok {
	//	// Если что-то пошло не так и user_id не смогли получить
	//	http.Error(w, "Не удалось получить user_id из контекста", http.StatusInternalServerError)
	//	return
	//}

	// Если всё ок, возвращаем сообщение, что доступ разрешён.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Доступ к защищённому маршруту разрешен.")))
}
