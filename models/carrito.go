package models

import (
	"context"
	"database/sql"
	"fmt"

	"tienda-online/internal/db"
	sqlutil "tienda-online/internal/sql"
)

type Carrito struct {
	IdCarrito int
	IdUsuario int
}

func (c Carrito) String() string {
	userLabel := fmt.Sprintf("UsuarioID:%d", c.IdUsuario)
	user, err := NewClienteManager(nil).Get(context.Background(), c.IdUsuario)
	if err == nil && user != nil {
		userLabel = user.String()
	}
	return fmt.Sprintf("[ Carrito #%d | %s ]", c.IdCarrito, userLabel)
}

type CarritoManager struct {
	db *sql.DB
}

func NewCarritoManager(database *sql.DB) *CarritoManager {
	if database == nil {
		database = db.CurrentDatabase
	}
	return &CarritoManager{db: database}
}

func (m *CarritoManager) List(ctx context.Context) ([]Carrito, error) {
	rows, err := db.QueryRowsFromFile(ctx, "leer/carrito.sql")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return sqlutil.ParseRow[Carrito](rows)
}

func (m *CarritoManager) Get(ctx context.Context, id int) (*Carrito, error) {
	if err := requirePositive("idCarrito", id); err != nil {
		return nil, err
	}
	rows, err := db.QueryRowsFromFile(ctx, "leer/carrito_por_id.sql", sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := sqlutil.ParseRow[Carrito](rows)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("carrito %d no encontrado", id)
	}
	return &items[0], nil
}

func (m *CarritoManager) Create(ctx context.Context, userId int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idUsuario", userId); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "añadir/carrito.sql", sql.Named("userId", userId))
	return err
}

func (m *CarritoManager) Update(ctx context.Context, id, userId int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idCarrito", id); err != nil {
		return err
	}
	if err := requirePositive("idUsuario", userId); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "editar/carrito.sql",
		sql.Named("id", id),
		sql.Named("userId", userId),
	)
	return err
}

func (m *CarritoManager) Delete(ctx context.Context, id int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idCarrito", id); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "remover/carrito.sql", sql.Named("id", id))
	return err
}
