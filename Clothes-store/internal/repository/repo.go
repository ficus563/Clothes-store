package repository

import (
	"context"
	"database/sql"
	"clothes-shop/internal/models"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// продукты
// добавить товар
func (r *Repository) CreateProduct(ctx context.Context, p models.Product) error {
	query := `INSERT INTO Product (ID, ProductName, ProductType, Gender, Color, ProductDescription, Brand, Season, ProductSize, Price) VALUES (?,?,?,?,?,?,?,?,?,?)`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.Type, p.Gender, p.Color, p.Description, p.Brand, p.Season, p.Size, p.Price)
	return err
}

// получить все товары (простая версия без сложных фильтров для экономии времени)
func (r *Repository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	query := `SELECT ID, ProductName, ProductType, Gender, Color, ProductDescription, Brand, Season, ProductSize, Price FROM Product`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Gender, &p.Color, &p.Description, &p.Brand, &p.Season, &p.Size, &p.Price)
		if err != nil { return nil, err }
		products = append(products, p)
	}
	return products, nil
}

// получить товар по id
func (r *Repository) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `SELECT ID, ProductName, ProductType, Gender, Color, ProductDescription, Brand, Season, ProductSize, Price FROM Product WHERE ID = ?`
	var p models.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Type, &p.Gender, &p.Color, &p.Description, &p.Brand, &p.Season, &p.Size, &p.Price)
	if err != nil {
		if err == sql.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &p, nil
}

// удалить товар
func (r *Repository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM Product WHERE ID = ?`, id)
	return err
}

// --- РОЛИ ---
// создать роль
func (r *Repository) CreateRole(ctx context.Context, role models.Role) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO Role (ID, Name) VALUES (?, ?)`, role.ID, role.Name)
	return err
}
// получить роли
func (r *Repository) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT ID, Name FROM Role`)
	if err != nil { return nil, err }
	defer rows.Close()
	var roles []models.Role
	for rows.Next() {
		var role models.Role
		rows.Scan(&role.ID, &role.Name)
		roles = append(roles, role)
	}
	return roles, nil
}

// корзина
// добавить в корзину
func (r *Repository) CreateCart(ctx context.Context, c models.Cart) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO Cart (ID, UserID, ProductID, Count) VALUES (?,?,?,?)`, c.ID, c.UserID, c.ProductID, c.Count)
	return err
}
// получить корзину юзера
func (r *Repository) GetCartByUserID(ctx context.Context, userID uuid.UUID) ([]models.Cart, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT ID, UserID, ProductID, Count FROM Cart WHERE UserID = ?`, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var cart []models.Cart
	for rows.Next() {
		var c models.Cart
		rows.Scan(&c.ID, &c.UserID, &c.ProductID, &c.Count)
		cart = append(cart, c)
	}
	return cart, nil
}
// удалить из корзины
func (r *Repository) DeleteCart(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM Cart WHERE ID = ?`, id)
	return err
}

// заказы
// создать заказ
func (r *Repository) CreateOrder(ctx context.Context, o models.Order) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO "Order" (ID, UserID, DateTimes, Status) VALUES (?,?,?,?)`, o.ID, o.UserID, o.DateTime, o.Status)
	return err
}
// получить все заказы
func (r *Repository) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT ID, UserID, DateTimes, Status FROM "Order"`)
	if err != nil { return nil, err }
	defer rows.Close()
	var orders []models.Order
	for rows.Next() {
		var o models.Order
		rows.Scan(&o.ID, &o.UserID, &o.DateTime, &o.Status)
		orders = append(orders, o)
	}
	return orders, nil
}