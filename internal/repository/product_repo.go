package repository

import (
	"context"
	"database/sql"

	"clothes-shop/internal/models"
	"github.com/google/uuid"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Метод создать
func (r *ProductRepository) Create(ctx context.Context, p models.Product) error {
	var query = `INSERT INTO Product(ID, ProductName, ProductType, Gender, Color, ProductDescription, Brand, Season, ProductSize, Price) 
	             VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := r.db.ExecContext(ctx, query, 
		p.ID, 
		p.Name, 
		p.Type, 
		p.Gender, 
		p.Color, 
		p.Description, 
		p.Brand, 
		p.Season, 
		p.Size, 
		p.Price,
	)
	return err
}

// Метод получить
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var query = `SELECT ID, ProductName, ProductType, Gender, Color, ProductDescription, Brand, Season, ProductSize, Price FROM Product WHERE ID = ?`
	
	var p models.Product
	row := r.db.QueryRowContext(ctx, query, id)
	
	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Type,
		&p.Gender,
		&p.Color,
		&p.Description,
		&p.Brand,
		&p.Season,
		&p.Size,
		&p.Price,
	)

	// Обработка ошибок
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
