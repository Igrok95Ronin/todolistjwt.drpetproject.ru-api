package routes

import (
	"fmt"
	"github.com/Igrok95Ronin/todolistjwt.drpetproject.ru-api.git/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

// Структура для таблицы all_notes
type AllNotes struct {
	ID        int64     `gorm:"primaryKey;column:id"` // Первичный ключ
	Note      string    `gorm:"column:note"`          // Поле заметки
	Completed bool      `gorm:"column:completed"`     // Статус выполнения
	UserID    int64     `gorm:"column:user_id"`       // Связь с таблицей users
	CreatedAt time.Time `gorm:"column:created_at"`    // Дата создания
}

// Структура для таблицы users
type Users struct {
	ID           int64     `gorm:"primaryKey;column:id"`    // Первичный ключ
	UserName     string    `gorm:"column:user_name;unique"` // Уникальное имя пользователя
	Email        string    `gorm:"column:email;unique"`     // Уникальный email
	PasswordHash string    `gorm:"column:password_hash"`    // Хеш пароля
	RefreshToken string    `gorm:"column:refresh_token"`    // Токен обновления (может быть NULL)
	CreatedAt    time.Time `gorm:"column:created_at"`       // Дата создания
}

func InitDB() (db *gorm.DB) {

	cfg := config.GetConfig()

	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d sslmode=%s client_encoding=UTF8", cfg.DB.User, cfg.DB.DBName, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.SslMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return nil
	}

	return db
}
