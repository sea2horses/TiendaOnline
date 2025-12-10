package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tienda-online/internal/db"
	sqlutil "tienda-online/internal/sql"
)

type Devolucion struct {
	IdDevolucion int
	IdPedido     int
	Fecha        time.Time
	Estado       sql.NullString
	Descripcion  sql.NullString
	Resolucion   sql.NullString
}

type DevolucionManager struct {
	db *sql.DB
}

func NewDevolucionManager(database *sql.DB) *DevolucionManager {
	if database == nil {
		database = db.CurrentDatabase
	}
	return &DevolucionManager{db: database}
}

func (m *DevolucionManager) List(ctx context.Context) ([]Devolucion, error) {
	rows, err := db.QueryRowsFromFile(ctx, "leer/devolucion.sql")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return sqlutil.ParseRow[Devolucion](rows)
}

func (m *DevolucionManager) Get(ctx context.Context, id int) (*Devolucion, error) {
	if err := requirePositive("idDevolucion", id); err != nil {
		return nil, err
	}
	rows, err := db.QueryRowsFromFile(ctx, "leer/devolucion_por_id.sql", sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := sqlutil.ParseRow[Devolucion](rows)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("devolucion %d no encontrada", id)
	}
	return &items[0], nil
}

func (m *DevolucionManager) Create(ctx context.Context, orderId int, fecha time.Time, estado, descripcion, resolucion string) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idPedido", orderId); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "añadir/devolucion.sql",
		sql.Named("orderId", orderId),
		sql.Named("date", fecha),
		sql.Named("status", optionalString(estado)),
		sql.Named("description", optionalString(descripcion)),
		sql.Named("resolution", optionalString(resolucion)),
	)
	return err
}

func (m *DevolucionManager) Update(ctx context.Context, id, orderId int, fecha time.Time, estado, descripcion, resolucion string) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idDevolucion", id); err != nil {
		return err
	}
	if err := requirePositive("idPedido", orderId); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "editar/devolucion.sql",
		sql.Named("id", id),
		sql.Named("orderId", orderId),
		sql.Named("date", fecha),
		sql.Named("status", optionalString(estado)),
		sql.Named("description", optionalString(descripcion)),
		sql.Named("resolution", optionalString(resolucion)),
	)
	return err
}

func (m *DevolucionManager) Delete(ctx context.Context, id int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idDevolucion", id); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "remover/devolucion.sql", sql.Named("id", id))
	return err
}
