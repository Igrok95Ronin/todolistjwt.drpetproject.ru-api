package models

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// Структура для таблицы all_notes
type AllNotes struct {
	ID        int64     `json:"ID" gorm:"primaryKey;column:id"`    // Первичный ключ
	Note      string    `json:"note" gorm:"column:note"`           // Поле заметки
	Completed bool      `json:"completed" gorm:"column:completed"` // Статус выполнения
	UserID    int64     `json:"userID" gorm:"column:user_id"`      // Связь с таблицей users
	CreatedAt time.Time `gorm:"column:created_at"`                 // Дата создания
}

// Структура для таблицы users
type Users struct {
	ID           int64     `json:"ID" gorm:"primaryKey;column:id"`           // Первичный ключ
	UserName     string    `json:"username" gorm:"column:user_name;unique"`  // Уникальное имя пользователя
	Email        string    `json:"email" gorm:"column:email;unique"`         // Уникальный email
	PasswordHash string    `json:"password" gorm:"column:password_hash"`     // Хеш пароля
	RefreshToken string    `json:"refreshToken" gorm:"column:refresh_token"` // Токен обновления (может быть NULL)
	CreatedAt    time.Time `gorm:"column:created_at"`                        // Дата создания
}

// MyClaims - своя структура для claim'ов JWT, включающая стандартные поля jwt.RegisteredClaims
// и ID пользователя (UserID), чтобы знать, кому принадлежит токен.
type MyClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}
