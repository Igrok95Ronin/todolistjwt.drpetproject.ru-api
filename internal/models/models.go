package models

import "time"

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
