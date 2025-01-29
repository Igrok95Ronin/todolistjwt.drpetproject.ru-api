package routes

import (
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/config"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/handlers"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
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
	router.GET("/", h.Home)
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
	//router.GET("/protected", AuthMiddleware(h.Protected))
	//
	//router.POST("/logout", h.Logout)
}
