package repository

import (
	"clothes-shop/internal/models"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// --- ТОВАРЫ ---
func (r *Repository) CreateProduct(ctx context.Context, p models.Product) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO Product (ID, ProductName, ProductType, Gender, Color, ProductDescription, Brand, Season, ProductSize, Price) VALUES (?,?,?,?,?,?,?,?,?,?)`, p.ID, p.Name, p.Type, p.Gender, p.Color, p.Description, p.Brand, p.Season, p.Size, p.Price)
	return err
}

func (r *Repository) GetProducts(ctx context.Context, pType, gender, brand string, limit, offset int, sortPrice string) ([]models.Product, error) {
	query := `SELECT ID, ProductName, ProductType, Gender, Color, ProductDescription, Brand, Season, ProductSize, Price FROM Product WHERE 1=1`
	args := []interface{}{}

	if pType != "" {
		query += " AND ProductType = ?"
		args = append(args, pType)
	}
	if gender != "" {
		query += " AND Gender = ?"
		args = append(args, gender)
	}
	if brand != "" {
		query += " AND Brand = ?"
		args = append(args, brand)
	}

	if sortPrice == "desc" {
		query += " ORDER BY Price DESC"
	} else {
		query += " ORDER BY Price ASC"
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		rows.Scan(&p.ID, &p.Name, &p.Type, &p.Gender, &p.Color, &p.Description, &p.Brand, &p.Season, &p.Size, &p.Price)
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *Repository) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRowContext(ctx, `SELECT ID, ProductName, ProductType, Gender, Color, ProductDescription, Brand, Season, ProductSize, Price FROM Product WHERE ID = ?`, id).Scan(&p.ID, &p.Name, &p.Type, &p.Gender, &p.Color, &p.Description, &p.Brand, &p.Season, &p.Size, &p.Price)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (r *Repository) UpdateProduct(ctx context.Context, p models.Product) error {
	_, err := r.db.ExecContext(ctx, `UPDATE Product SET ProductName=?, ProductType=?, Gender=?, Color=?, ProductDescription=?, Brand=?, Season=?, ProductSize=?, Price=? WHERE ID=?`, p.Name, p.Type, p.Gender, p.Color, p.Description, p.Brand, p.Season, p.Size, p.Price, p.ID)
	return err
}

func (r *Repository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM Product WHERE ID = ?`, id)
	return err
}

// --- ФОТО ---
func (r *Repository) GetPhotoByID(ctx context.Context, id uuid.UUID) (*models.Photo, error) {
	var p models.Photo
	err := r.db.QueryRowContext(ctx, `SELECT ID, ProductID, Filename FROM Photo WHERE ID = ?`, id).Scan(&p.ID, &p.ProductID, &p.Filename)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (r *Repository) GetPhotoByFilename(ctx context.Context, filename string) (*models.Photo, error) {
	var p models.Photo
	err := r.db.QueryRowContext(ctx, `SELECT ID, ProductID, Filename FROM Photo WHERE Filename = ?`, filename).Scan(&p.ID, &p.ProductID, &p.Filename)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (r *Repository) CreatePhoto(ctx context.Context, p models.Photo) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO Photo (ID, ProductID, Filename) VALUES (?,?,?)`, p.ID, p.ProductID, p.Filename)
	return err
}

func (r *Repository) DeletePhoto(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM Photo WHERE ID = ?`, id)
	return err
}

// --- ПОЛЬЗОВАТЕЛИ И РОЛИ ---
func (r *Repository) CreateUser(ctx context.Context, u models.User) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO "User" (ID, NameUser, Email, Phone, UserRole, RegistrationDate, PasswordHash) VALUES (?,?,?,?,?,?,?)`, u.ID, u.Name, u.Email, u.Phone, u.RoleID, u.RegistrationDate, u.Password)
	return err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, string, error) {
	var u models.User
	var roleName string
	query := `SELECT u.ID, u.NameUser, u.Email, u.UserRole, u.PasswordHash, r.Name FROM "User" u JOIN Role r ON u.UserRole = r.ID WHERE u.Email = ?`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Name, &u.Email, &u.RoleID, &u.Password, &roleName)
	if err == sql.ErrNoRows {
		return nil, "", nil
	}
	return &u, roleName, err
}

func (r *Repository) GetUsers(ctx context.Context, limit, offset int) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT ID, NameUser, Email, Phone, UserRole, RegistrationDate FROM "User" LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var us []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.RoleID, &u.RegistrationDate)
		us = append(us, u)
	}
	return us, rows.Err()
}

func (r *Repository) GetUsersByRole(ctx context.Context, roleID uuid.UUID) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT ID, NameUser, Email FROM "User" WHERE UserRole = ?`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var us []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Name, &u.Email)
		us = append(us, u)
	}
	return us, rows.Err()
}

func (r *Repository) UpdateUser(ctx context.Context, u models.User) error {
	_, err := r.db.ExecContext(ctx, `UPDATE "User" SET NameUser=?, Email=?, Phone=?, UserRole=? WHERE ID=?`, u.Name, u.Email, u.Phone, u.RoleID, u.ID)
	return err
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM "User" WHERE ID = ?`, id)
	return err
}

func (r *Repository) CreateRole(ctx context.Context, role models.Role) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO Role (ID, Name) VALUES (?, ?)`, role.ID, role.Name)
	return err
}

func (r *Repository) GetRoles(ctx context.Context) ([]models.Role, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT ID, Name FROM Role`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rs []models.Role
	for rows.Next() {
		var role models.Role
		rows.Scan(&role.ID, &role.Name)
		rs = append(rs, role)
	}
	return rs, rows.Err()
}

func (r *Repository) UpdateRole(ctx context.Context, role models.Role) error {
	_, err := r.db.ExecContext(ctx, `UPDATE Role SET Name = ? WHERE ID = ?`, role.Name, role.ID)
	return err
}

func (r *Repository) DeleteRole(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM Role WHERE ID = ?`, id)
	return err
}

// --- КОРЗИНА ---
func (r *Repository) CreateCart(ctx context.Context, c models.Cart) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO Cart (ID, UserID, ProductID, Count) VALUES (?,?,?,?)`, c.ID, c.UserID, c.ProductID, c.Count)
	return err
}

func (r *Repository) GetCartByUserID(ctx context.Context, userID uuid.UUID) ([]models.CartResponse, error) {
	query := `SELECT c.ID, p.ProductName, p.Price, c.Count FROM Cart c JOIN Product p ON c.ProductID = p.ID WHERE c.UserID = ?`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cart []models.CartResponse
	for rows.Next() {
		var c models.CartResponse
		rows.Scan(&c.ID, &c.ProductName, &c.Price, &c.Count)
		cart = append(cart, c)
	}
	return cart, rows.Err()
}

func (r *Repository) UpdateCart(ctx context.Context, c models.Cart) error {
	_, err := r.db.ExecContext(ctx, `UPDATE Cart SET Count = ? WHERE ID = ?`, c.Count, c.ID)
	return err
}

func (r *Repository) DeleteCart(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM Cart WHERE ID = ?`, id)
	return err
}

// --- ЗАКАЗЫ ---
func (r *Repository) CreateOrder(ctx context.Context, o models.Order) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO "Order" (ID, UserID, DateTimes, Status) VALUES (?,?,?,?)`, o.ID, o.UserID, o.DateTime, o.Status)
	return err
}

func (r *Repository) GetOrders(ctx context.Context, limit, offset int) ([]models.Order, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT ID, UserID, DateTimes, Status FROM "Order" ORDER BY DateTimes DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var os []models.Order
	for rows.Next() {
		var o models.Order
		rows.Scan(&o.ID, &o.UserID, &o.DateTime, &o.Status)
		os = append(os, o)
	}
	return os, rows.Err()
}

func (r *Repository) GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var o models.Order
	err := r.db.QueryRowContext(ctx, `SELECT ID, UserID, DateTimes, Status FROM "Order" WHERE ID = ?`, id).Scan(&o.ID, &o.UserID, &o.DateTime, &o.Status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &o, err
}

func (r *Repository) UpdateOrder(ctx context.Context, o models.Order) error {
	_, err := r.db.ExecContext(ctx, `UPDATE "Order" SET Status = ? WHERE ID = ?`, o.Status, o.ID)
	return err
}

// --- ЭЛЕМЕНТЫ ЗАКАЗА ---
func (r *Repository) CreateOrderItem(ctx context.Context, item models.OrderItem) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO OrderItem (ID, OrderID, ProductID, Price) VALUES (?,?,?,?)`, item.ID, item.OrderID, item.ProductID, item.Price)
	return err
}

func (r *Repository) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]models.OrderItem, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT ID, OrderID, ProductID, Price FROM OrderItem WHERE OrderID = ?`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.OrderItem
	for rows.Next() {
		var i models.OrderItem
		rows.Scan(&i.ID, &i.OrderID, &i.ProductID, &i.Price)
		items = append(items, i)
	}
	return items, rows.Err()
}
