package models

import (
	"context"
	"database/sql"
	"fmt"

	"tienda-online/internal/db"
	sqlutil "tienda-online/internal/sql"
)

type Categoria struct {
	IdCategoria int
	Nombre      string
}

func (c Categoria) String() string {
	return fmt.Sprintf("Categoria #%d | %s", c.IdCategoria, c.Nombre)
}

type CategoriaManager struct {
	db *sql.DB
}

func NewCategoriaManager(database *sql.DB) *CategoriaManager {
	if database == nil {
		database = db.CurrentDatabase
	}
	return &CategoriaManager{db: database}
}

func (m *CategoriaManager) List(ctx context.Context) ([]Categoria, error) {
	rows, err := db.QueryRowsFromFile(ctx, "leer/categoria.sql")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return sqlutil.ParseRow[Categoria](rows)
}

func (m *CategoriaManager) Get(ctx context.Context, id int) (*Categoria, error) {
	if err := requirePositive("idCategoria", id); err != nil {
		return nil, err
	}
	rows, err := db.QueryRowsFromFile(ctx, "leer/categoria_por_id.sql", sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := sqlutil.ParseRow[Categoria](rows)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("categoria %d no encontrada", id)
	}
	return &items[0], nil
}

func (m *CategoriaManager) Create(ctx context.Context, name string) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	name, err := requireNonEmpty("nombre", name)
	if err != nil {
		return err
	}
	_, err = db.ExecFromFile(ctx, "añadir/categoria.sql", sql.Named("name", name))
	return err
}

func (m *CategoriaManager) Update(ctx context.Context, id int, name string) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idCategoria", id); err != nil {
		return err
	}
	name, err := requireNonEmpty("nombre", name)
	if err != nil {
		return err
	}
	_, err = db.ExecFromFile(ctx, "editar/categoria.sql",
		sql.Named("id", id),
		sql.Named("name", name),
	)
	return err
}

func (m *CategoriaManager) Delete(ctx context.Context, id int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idCategoria", id); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "remover/categoria.sql", sql.Named("id", id))
	return err
}
