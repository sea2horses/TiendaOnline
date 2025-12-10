package models

import (
	"context"
	"database/sql"
	"fmt"

	"tienda-online/internal/db"
	sqlutil "tienda-online/internal/sql"
)

type Pedido struct {
	IdPedido  int
	IdUsuario int
	Entregado bool
}

func (p Pedido) String() string {
	ent := "No entregado"
	if p.Entregado {
		ent = "Entregado"
	}
	return fmt.Sprintf("Pedido #%d | UsuarioID:%d | %s", p.IdPedido, p.IdUsuario, ent)
}

type PedidoManager struct {
	db *sql.DB
}

func NewPedidoManager(database *sql.DB) *PedidoManager {
	if database == nil {
		database = db.CurrentDatabase
	}
	return &PedidoManager{db: database}
}

func (m *PedidoManager) List(ctx context.Context) ([]Pedido, error) {
	rows, err := db.QueryRowsFromFile(ctx, "leer/pedido.sql")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return sqlutil.ParseRow[Pedido](rows)
}

func (m *PedidoManager) Get(ctx context.Context, id int) (*Pedido, error) {
	if err := requirePositive("idPedido", id); err != nil {
		return nil, err
	}
	rows, err := db.QueryRowsFromFile(ctx, "leer/pedido_por_id.sql", sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := sqlutil.ParseRow[Pedido](rows)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("pedido %d no encontrado", id)
	}
	return &items[0], nil
}

func (m *PedidoManager) Create(ctx context.Context, userId int, delivered bool) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idUsuario", userId); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "añadir/pedido.sql",
		sql.Named("userId", userId),
		sql.Named("delivered", delivered),
	)
	return err
}

func (m *PedidoManager) Update(ctx context.Context, id, userId int, delivered bool) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idPedido", id); err != nil {
		return err
	}
	if err := requirePositive("idUsuario", userId); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "editar/pedido.sql",
		sql.Named("id", id),
		sql.Named("userId", userId),
		sql.Named("delivered", delivered),
	)
	return err
}

func (m *PedidoManager) Delete(ctx context.Context, id int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idPedido", id); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "remover/pedido.sql", sql.Named("id", id))
	return err
}
