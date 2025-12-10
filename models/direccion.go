package models

import (
	"context"
	"database/sql"
	"fmt"


	"tienda-online/internal"
	"tienda-online/internal/db"
	sqlutil "tienda-online/internal/sql"
)

type Direccion struct {
	IdDirección int
	IdUsuario   int
	Tipo        string
	Detalle     sql.NullString
}

func (d Direccion) String() string {
	return fmt.Sprintf("Direccion #%d | UsuarioID:%d | %s | %s", d.IdDirección, d.IdUsuario, d.Tipo, internal.NullString(d.Detalle))
}

type DireccionManager struct {
	db *sql.DB
}

func NewDireccionManager(database *sql.DB) *DireccionManager {
	if database == nil {
		database = db.CurrentDatabase
	}
	return &DireccionManager{db: database}
}

func (m *DireccionManager) List(ctx context.Context) ([]Direccion, error) {
	rows, err := db.QueryRowsFromFile(ctx, "leer/direccion.sql")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return sqlutil.ParseRow[Direccion](rows)
}

func (m *DireccionManager) Get(ctx context.Context, id int) (*Direccion, error) {
	if err := requirePositive("idDireccion", id); err != nil {
		return nil, err
	}
	rows, err := db.QueryRowsFromFile(ctx, "leer/direccion_por_id.sql", sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := sqlutil.ParseRow[Direccion](rows)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("direccion %d no encontrada", id)
	}
	return &items[0], nil
}

func (m *DireccionManager) Create(ctx context.Context, userId int, tipo, detalle string) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idUsuario", userId); err != nil {
		return err
	}
	tipo, err := requireNonEmpty("tipo", tipo)
	if err != nil {
		return err
	}
	_, err = db.ExecFromFile(ctx, "añadir/direccion.sql",
		sql.Named("userId", userId),
		sql.Named("type", tipo),
		sql.Named("detail", optionalString(detalle)),
	)
	return err
}

func (m *DireccionManager) Update(ctx context.Context, id, userId int, tipo, detalle string) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idDireccion", id); err != nil {
		return err
	}
	if err := requirePositive("idUsuario", userId); err != nil {
		return err
	}
	tipo, err := requireNonEmpty("tipo", tipo)
	if err != nil {
		return err
	}
	_, err = db.ExecFromFile(ctx, "editar/direccion.sql",
		sql.Named("id", id),
		sql.Named("userId", userId),
		sql.Named("type", tipo),
		sql.Named("detail", optionalString(detalle)),
	)
	return err
}

func (m *DireccionManager) Delete(ctx context.Context, id int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idDireccion", id); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "remover/direccion.sql", sql.Named("id", id))
	return err
}
