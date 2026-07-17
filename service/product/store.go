package product

import (
	"database/sql"

	"github.com/justKody/ecom-go/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []types.Product{}
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	return products, nil
}

func (s *Store) CreateProduct(product types.CreateProductPayload) error {
	_, err := s.db.Exec("INSERT INTO products (name, description, imageUrl, price, quantity) VALUES (?, ?, ?, ?, ?)", product.Name, product.Description, product.ImageURL, product.Price, product.Quantity)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetProductByIDs(ids []int) ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []types.Product{}
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	return products, nil
}

func (s *Store) ReduceProductQuantity(id int, quantity int) error {
	_, err := s.db.Exec("UPDATE products SET quantity = quantity - ? WHERE id = ?", quantity, id)
	if err != nil {
		return err
	}
	return nil
}
func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	var product types.Product
	err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.ImageURL, &product.Price, &product.Quantity, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
