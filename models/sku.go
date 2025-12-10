package models

import (
	"context"
	"database/sql"
	"fmt"

	"tienda-online/internal/db"
	sqlutil "tienda-online/internal/sql"
)

type SKU struct {
	IdSKU      int
	IdProducto int
	Precio     float64
	Stock      int
}

type SKUManager struct {
	db *sql.DB
}

func NewSKUManager(database *sql.DB) *SKUManager {
	if database == nil {
		database = db.CurrentDatabase
	}
	return &SKUManager{db: database}
}

func (m *SKUManager) List(ctx context.Context) ([]SKU, error) {
	rows, err := db.QueryRowsFromFile(ctx, "leer/sku.sql")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return sqlutil.ParseRow[SKU](rows)
}

func (m *SKUManager) Get(ctx context.Context, id int) (*SKU, error) {
	if err := requirePositive("idSKU", id); err != nil {
		return nil, err
	}
	rows, err := db.QueryRowsFromFile(ctx, "leer/sku_por_id.sql", sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := sqlutil.ParseRow[SKU](rows)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("SKU %d no encontrado", id)
	}
	return &items[0], nil
}

func (m *SKUManager) Create(ctx context.Context, productId int, price float64, stock int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idProducto", productId); err != nil {
		return err
	}
	if price < 0 {
		return fmt.Errorf("precio no puede ser negativo")
	}
	if stock < 0 {
		return fmt.Errorf("stock no puede ser negativo")
	}
	_, err := db.ExecFromFile(ctx, "añadir/sku.sql",
		sql.Named("productId", productId),
		sql.Named("price", price),
		sql.Named("stock", stock),
	)
	return err
}

func (m *SKUManager) Update(ctx context.Context, id, productId int, price float64, stock int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idSKU", id); err != nil {
		return err
	}
	if err := requirePositive("idProducto", productId); err != nil {
		return err
	}
	if price < 0 {
		return fmt.Errorf("precio no puede ser negativo")
	}
	if stock < 0 {
		return fmt.Errorf("stock no puede ser negativo")
	}
	_, err := db.ExecFromFile(ctx, "editar/sku.sql",
		sql.Named("id", id),
		sql.Named("productId", productId),
		sql.Named("price", price),
		sql.Named("stock", stock),
	)
	return err
}

func (m *SKUManager) Delete(ctx context.Context, id int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idSKU", id); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "remover/sku.sql", sql.Named("id", id))
	return err
}
