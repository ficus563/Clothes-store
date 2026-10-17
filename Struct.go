package	main
import {
	"github.com/google/uuid"
}

//Структура товара
type Product struct {
	ID int
	Name, Type, Gender, Color, Description, Brand, Season string
	Size float64
	Price float64
}

//Структура фото
type Photo struct {
	ID, ProductID int
	Filename string
}

//Структура роль пользователя
type Role struct {
	ID int
	Name string
}

//Структура пользователь
type User struct {
	ID int
	Name, Email, Phone, Role string
	RegistrationDate string // выбрать корректный тип данных
}

//Структура корзина
type Cart struct {
	ID int 
	UserID, ProductID int
	Count float64
}

//Структура заказ
type Order struct {
	ID int
	UserID int
	DateTime string
	Status string
}

//Cтруктура элементы заказа
type  OrderItem struct {
	ID int
	OrderID, ProductID int
	Price float64
}