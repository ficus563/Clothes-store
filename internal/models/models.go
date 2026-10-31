package models

import (
	"time"
	"github.com/google/uuid"
)

// Структура товара
type Product struct {
	ID          uuid.UUID
	Name        string
	Type        string
	Gender      string
	Color       string
	Description string
	Brand       string
	Season      string
	Size        int
	Price       float64
}

// Структура фото
type Photo struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Filename  string
}

// Структура роль пользователя
type Role struct {
	ID   uuid.UUID
	Name string
}

// Структура пользователь
type User struct {
	ID               uuid.UUID
	Name             string
	Email            string
	Phone            string
	RoleID           uuid.UUID // Ссылка на роль
	RegistrationDate time.Time
}

// Структура корзина
type Cart struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ProductID uuid.UUID
	Count     int
}

// Статусы заказов
type OrderStatus string

const (
	StatusAccepted   OrderStatus = "Принят"
	StatusInProgress OrderStatus = "В работе"
	StatusAssembled  OrderStatus = "Собран"
	StatusReady      OrderStatus = "Готов к выдаче"
	StatusCancelled  OrderStatus = "Отменён"
)

// Структура заказ
type Order struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	DateTime time.Time
	Status   OrderStatus
}

// Структура элементы заказа
type OrderItem struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	ProductID uuid.UUID
	Price     float64
}
