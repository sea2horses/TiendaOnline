package models

import (
	"context"
	"database/sql"
	"fmt"

	"tienda-online/internal"
	"tienda-online/internal/db"
	sqlutil "tienda-online/internal/sql"
)

type Resena struct {
	IdReseña   int
	IdUsuario  int
	IdProducto int
	Puntuación int
	Comentario sql.NullString
}

func (r Resena) String() string {
	return fmt.Sprintf("[ Reseña #%d | UsuarioID:%d | ProductoID:%d | %d/5 | %s ]",
		r.IdReseña, r.IdUsuario, r.IdProducto, r.Puntuación, internal.NullString(r.Comentario))
}

type ResenaManager struct {
	db *sql.DB
}

func NewResenaManager(database *sql.DB) *ResenaManager {
	if database == nil {
		database = db.CurrentDatabase
	}
	return &ResenaManager{db: database}
}

func (m *ResenaManager) List(ctx context.Context) ([]Resena, error) {
	rows, err := db.QueryRowsFromFile(ctx, "leer/resena.sql")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return sqlutil.ParseRow[Resena](rows)
}

func (m *ResenaManager) Get(ctx context.Context, id int) (*Resena, error) {
	if err := requirePositive("idReseña", id); err != nil {
		return nil, err
	}
	rows, err := db.QueryRowsFromFile(ctx, "leer/resena_por_id.sql", sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := sqlutil.ParseRow[Resena](rows)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("reseña %d no encontrada", id)
	}
	return &items[0], nil
}

func (m *ResenaManager) Create(ctx context.Context, userId, productId, rating int, comment string) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idUsuario", userId); err != nil {
		return err
	}
	if err := requirePositive("idProducto", productId); err != nil {
		return err
	}
	if rating < 1 || rating > 5 {
		return fmt.Errorf("puntuacion debe estar entre 1 y 5")
	}
	_, err := db.ExecFromFile(ctx, "añadir/resena.sql",
		sql.Named("userId", userId),
		sql.Named("productId", productId),
		sql.Named("rating", rating),
		sql.Named("comment", optionalString(comment)),
	)
	return err
}

func (m *ResenaManager) Update(ctx context.Context, id, userId, productId, rating int, comment string) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idReseña", id); err != nil {
		return err
	}
	if err := requirePositive("idUsuario", userId); err != nil {
		return err
	}
	if err := requirePositive("idProducto", productId); err != nil {
		return err
	}
	if rating < 1 || rating > 5 {
		return fmt.Errorf("puntuacion debe estar entre 1 y 5")
	}
	_, err := db.ExecFromFile(ctx, "editar/resena.sql",
		sql.Named("id", id),
		sql.Named("userId", userId),
		sql.Named("productId", productId),
		sql.Named("rating", rating),
		sql.Named("comment", optionalString(comment)),
	)
	return err
}

func (m *ResenaManager) Delete(ctx context.Context, id int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idReseña", id); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "remover/resena.sql", sql.Named("id", id))
	return err
}
