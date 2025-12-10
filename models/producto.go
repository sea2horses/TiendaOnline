package models

import (
	"context"
	"database/sql"
	"fmt"

	"tienda-online/internal/db"
	sqlutil "tienda-online/internal/sql"
)

type Producto struct {
	IdProducto  int
	Descripcion string
	IdCategoria sql.NullInt32
}

func (p Producto) String() string {
	cat := "sin categoria"
	if p.IdCategoria.Valid {
		cat = fmt.Sprintf("CategoriaID:%d", p.IdCategoria.Int32)
	}
	return fmt.Sprintf("[ Producto #%d | %s | %s ]", p.IdProducto, p.Descripcion, cat)
}

type ProductoManager struct {
	db *sql.DB
}

func NewProductoManager(database *sql.DB) *ProductoManager {
	if database == nil {
		database = db.CurrentDatabase
	}
	return &ProductoManager{db: database}
}

func (m *ProductoManager) List(ctx context.Context) ([]Producto, error) {
	rows, err := db.QueryRowsFromFile(ctx, "leer/producto.sql")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return sqlutil.ParseRow[Producto](rows)
}

func (m *ProductoManager) Get(ctx context.Context, id int) (*Producto, error) {
	if err := requirePositive("idProducto", id); err != nil {
		return nil, err
	}
	rows, err := db.QueryRowsFromFile(ctx, "leer/producto_por_id.sql", sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := sqlutil.ParseRow[Producto](rows)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("producto %d no encontrado", id)
	}
	return &items[0], nil
}

func (m *ProductoManager) Create(ctx context.Context, description string, categoryId int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	description, err := requireNonEmpty("descripcion", description)
	if err != nil {
		return err
	}
	_, err = db.ExecFromFile(ctx, "añadir/producto.sql",
		sql.Named("description", description),
		sql.Named("categoryId", optionalInt(categoryId)),
	)
	return err
}

func (m *ProductoManager) Update(ctx context.Context, id int, description string, categoryId int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idProducto", id); err != nil {
		return err
	}
	description, err := requireNonEmpty("descripcion", description)
	if err != nil {
		return err
	}
	_, err = db.ExecFromFile(ctx, "editar/producto.sql",
		sql.Named("id", id),
		sql.Named("description", description),
		sql.Named("categoryId", optionalInt(categoryId)),
	)
	return err
}

func (m *ProductoManager) Delete(ctx context.Context, id int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idProducto", id); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "remover/producto.sql", sql.Named("id", id))
	return err
}
