package	main
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
	ProductID int
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
	Name, Email, Phone, Role string
	RegistrationDate time.Time // выбрать корректный тип данных
}

//Структура корзина
type Cart struct {
	ID uuid.UUID
	UserID, ProductID, Count int

}

//Структура заказ
type Order struct {
	ID uuid.UUID
	UserID int
	DateTime time.Time
	Status string
}

//Cтруктура элементы заказа
type  OrderItem struct {
	ID uuid.UUID
	OrderID, ProductID int
	Price float64
}