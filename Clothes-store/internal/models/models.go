package models

import (
	"time"

	"github.com/google/uuid"
)

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

type Role struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type User struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	RoleID           uuid.UUID `json:"role_id"`
	RegistrationDate time.Time `json:"registration_date"`
	Password         string    `json:"password,omitempty"`
}

// Для входа и проверки ролей
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	RoleName string    `json:"role_name"`
}

type Cart struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	Count     int       `json:"count"`
}
type CartResponse struct {
	ID          uuid.UUID `json:"id"`
	ProductName string    `json:"product_name"`
	Price       float64   `json:"price"`
	Count       int       `json:"count"`
}

type OrderStatus string

const (
	StatusAccepted   OrderStatus = "Принят"
	StatusInProgress OrderStatus = "В работе"
	StatusAssembled  OrderStatus = "Собран"
	StatusReady      OrderStatus = "Готов к выдаче"
	StatusCancelled  OrderStatus = "Отменён"
)

type Order struct {
	ID       uuid.UUID   `json:"id"`
	UserID   uuid.UUID   `json:"user_id"`
	DateTime time.Time   `json:"date_time"`
	Status   OrderStatus `json:"status"`
}
type OrderItem struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Price     float64   `json:"price"`
}
type Photo struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Filename  string    `json:"filename"`
}
