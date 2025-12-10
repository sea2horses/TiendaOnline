package models

import (
	"context"
	"database/sql"
	"fmt"

	"tienda-online/internal"
	"tienda-online/internal/db"
	sqlutil "tienda-online/internal/sql"
)

// Cliente representa la fila de la tabla Clientes.
type Cliente struct {
	IdUsuario    int
	Nombre       string
	Telefono     string
	Correo       sql.NullString
	PasswordHash []byte
	PasswordSalt []byte
}

func (c Cliente) String() string {
	return fmt.Sprintf("id: %d, nombre: %s, telefono: %s, correo: %v.", c.IdUsuario, c.Nombre, c.Telefono, internal.NullString(c.Correo))
}

type ClienteManager struct {
	db *sql.DB
}

func NewClienteManager(database *sql.DB) *ClienteManager {
	if database == nil {
		database = db.CurrentDatabase
	}
	return &ClienteManager{db: database}
}

// List obtiene todos los clientes.
func (m *ClienteManager) List(ctx context.Context) ([]Cliente, error) {
	rows, err := db.QueryRowsFromFile(ctx, "leer/cliente.sql")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return sqlutil.ParseRow[Cliente](rows)
}

// Get trae un cliente por ID.
func (m *ClienteManager) Get(ctx context.Context, id int) (*Cliente, error) {
	if err := requirePositive("idUsuario", id); err != nil {
		return nil, err
	}
	rows, err := db.QueryRowsFromFile(ctx, "leer/cliente_por_id.sql", sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := sqlutil.ParseRow[Cliente](rows)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("cliente %d no encontrado", id)
	}
	return &items[0], nil
}

// Create inserta un nuevo cliente (hash y salt se generan en la query).
func (m *ClienteManager) Create(ctx context.Context, name, phone, email, password string) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	var err error
	name, err = requireNonEmpty("nombre", name)
	if err != nil {
		return err
	}
	phone, err = requireNonEmpty("telefono", phone)
	if err != nil {
		return err
	}
	password, err = requireNonEmpty("contraseña", password)
	if err != nil {
		return err
	}

	_, err = db.ExecFromFile(ctx, "añadir/cliente.sql",
		sql.Named("name", name),
		sql.Named("phone", phone),
		sql.Named("email", optionalString(email)),
		sql.Named("password", password),
	)
	return err
}

// Update modifica datos básicos del cliente (sin contraseña).
func (m *ClienteManager) Update(ctx context.Context, id int, name, phone, email string) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idUsuario", id); err != nil {
		return err
	}
	var err error
	name, err = requireNonEmpty("nombre", name)
	if err != nil {
		return err
	}
	phone, err = requireNonEmpty("telefono", phone)
	if err != nil {
		return err
	}
	_, err = db.ExecFromFile(ctx, "editar/cliente.sql",
		sql.Named("id", id),
		sql.Named("name", name),
		sql.Named("phone", phone),
		sql.Named("email", optionalString(email)),
	)
	return err
}

// UpdatePassword cambia la contraseña aplicando hash/salt en el SQL.
func (m *ClienteManager) UpdatePassword(ctx context.Context, id int, password string) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idUsuario", id); err != nil {
		return err
	}
	var err error
	password, err = requireNonEmpty("contraseña", password)
	if err != nil {
		return err
	}

	_, err = db.ExecFromFile(ctx, "editar/cliente_contraseña.sql",
		sql.Named("id", id),
		sql.Named("password", password),
	)
	return err
}

// Delete borra un cliente por ID.
func (m *ClienteManager) Delete(ctx context.Context, id int) error {
	if err := ensureDB(m.db); err != nil {
		return err
	}
	if err := requirePositive("idUsuario", id); err != nil {
		return err
	}
	_, err := db.ExecFromFile(ctx, "remover/cliente.sql", sql.Named("id", id))
	return err
}
