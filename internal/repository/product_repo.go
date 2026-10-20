package repository
import (
	"context"
	"database-sql"
	"clothes-shop/internal/models"
)

type ProductRepository struct {
	db *sql.DB 
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

//Метод создать
func (r *ProductRepository) Create(ctx context.Context, p models.Product) error {
	var query = 'INSERT INTO Product(ID, ProductName, ProductType, Gender, Color, ProductDescription, Brand, Season, ProductSize, Price) VALUES ('?','?','?','?','?','?','?','?','?','?')'
	_, err := r.db.ExecContext(ctx, query, p.ID, p.ProductName, p.ProductType, p.Gender, p.Color, p.ProductDescription, p.Brand, p.Season, p.ProductSize, p.Price)
	return err 
}

//func (g *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	//var p models.Product
	//var row r.db.QueryRowContext(ctx, query, id)
	//var query = 'SELECT ID, ProductName, ProductType, Gender, Color, ProductDescription, Brand, Season, ProductSize, Price FROM Product WHERE ID = ?'


//}

//TODO
