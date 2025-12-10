package models

import (
	"context"
	"database/sql"
	"fmt"

	"tienda-online/internal/db"
	sqlutil "tienda-online/internal/sql"
)

type CarritoDetalle struct {
	IdDetalle int
	IdCarrito int
	IdSKU     int
	Cantidad  int
}

type CarritoDetalleManager struct {
	db *sql.DB
}

func NewCarritoDetalleManager(database *sql.DB) *CarritoDetalleManager {
	if database == nil {
		database = db.CurrentDatabase
	}
	return &CarritoDetalleManager{db: database}
}

func (m *CarritoDetalleManager) List(ctx context.Context) ([]CarritoDetalle, error) {
	rows, err := db.QueryRowsFromFile(ctx, "leer/carrito_detalle.sql")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return sqlutil.ParseRow[CarritoDetalle](rows)
}

func (m *CarritoDetalleManager) Get(ctx context.Context, id int) (*CarritoDetalle, error) {
	if err := requirePositive("idDetalle", id); err != nil {
		return nil, err
	}
	rows, err := db.QueryRowsFromFile(ctx, "leer/carrito_detalle_por_id.sql", sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := sqlutil.ParseRow[CarritoDetalle](rows)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("detalle %d no encontrado", id)
	}
	return &items[0], nil
}

func (m *CarritoDetalleManager) Create(ctx context.Context, cartId, skuId, quantity int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idCarrito", cartId); err != nil {
		return err
	}
	if err := requirePositive("idSKU", skuId); err != nil {
		return err
	}
	if quantity <= 0 {
		return fmt.Errorf("cantidad debe ser mayor a cero")
	}
	_, err := db.ExecFromFile(ctx, "añadir/carrito_detalle.sql",
		sql.Named("cartId", cartId),
		sql.Named("skuId", skuId),
		sql.Named("quantity", quantity),
	)
	return err
}

func (m *CarritoDetalleManager) Update(ctx context.Context, id, cartId, skuId, quantity int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idDetalle", id); err != nil {
		return err
	}
	if err := requirePositive("idCarrito", cartId); err != nil {
		return err
	}
	if err := requirePositive("idSKU", skuId); err != nil {
		return err
	}
	if quantity <= 0 {
		return fmt.Errorf("cantidad debe ser mayor a cero")
	}
	_, err := db.ExecFromFile(ctx, "editar/carrito_detalle.sql",
		sql.Named("id", id),
		sql.Named("cartId", cartId),
		sql.Named("skuId", skuId),
		sql.Named("quantity", quantity),
	)
	return err
}

func (m *CarritoDetalleManager) Delete(ctx context.Context, id int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idDetalle", id); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "remover/carrito_detalle.sql", sql.Named("id", id))
	return err
}
