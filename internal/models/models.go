package	models
import {
	"github.com/google/uuid"
	"time"
}

//Структура товара
type Product struct {
	ID uuid.UUID
	Name, Type, Gender, Color, Description, Brand, Season string
	Size int
	Price float64
}

//Структура фото
type Photo struct {
	ID uuid.UUID
	ProductID uuid.UUID
	Filename string
}

//Структура роль пользователя
type Role struct {
	ID uuid.UUID
	Name string
}

//Структура пользователь
type User struct {
	ID uuid.UUID
	Name, Email, Phone string
	Role uuid.UUID
	RegistrationDate time.Time
}

//Структура корзина
type Cart struct {
	ID uuid.UUID
	UserID, ProductID uuid.UUID
	Count int

}

//Статусы заказов
type OrderStatus string 
const (
	StatusAccepted OrderStatus = "Принят"
	StatusInProgress OrderStatus = "В работе"
	StatusAssembled OrderStatus = "Собран"
	StatusReady OrderStatus = "Готов к выдаче"
	StatusCancelled OrderStatus = "Отменён"
)


//Структура заказ
type Order struct {
	ID uuid.UUID
	UserID uuid.UUID
	DateTime time.Time
	Status OrderStatus
}

//Cтруктура элементы заказа
type  OrderItem struct {
	ID uuid.UUID
	OrderID, ProductID uuid.UUID
	Price float64
}
