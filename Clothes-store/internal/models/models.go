package models

import (
	"time"
	"github.com/google/uuid"
)

// структура товара
type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Gender      string    `json:"gender"`
	Color       string    `json:"color"`
	Description string    `json:"description"`
	Brand       string    `json:"brand"`
	Season      string    `json:"season"`
	Size        int       `json:"size"`
	Price       float64   `json:"price"`
}

// структура фото
type Photo struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Filename  string    `json:"filename"`
}

// структура роль
type Role struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// структура пользователь
type User struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	RoleID           uuid.UUID `json:"role_id"`
	RegistrationDate time.Time `json:"registration_date"`
}

// структура корзина
type Cart struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	Count     int       `json:"count"`
}

// статусы заказа
type OrderStatus string
const (
	StatusAccepted   OrderStatus = "Принят"
	StatusInProgress OrderStatus = "В работе"
	StatusAssembled  OrderStatus = "Собран"
	StatusReady      OrderStatus = "Готов к выдаче"
	StatusCancelled  OrderStatus = "Отменён"
)

// структура заказ
type Order struct {
	ID       uuid.UUID   `json:"id"`
	UserID   uuid.UUID   `json:"user_id"`
	DateTime time.Time   `json:"date_time"`
	Status   OrderStatus `json:"status"`
}

// структура элемент заказа
type OrderItem struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Price     float64   `json:"price"`
}