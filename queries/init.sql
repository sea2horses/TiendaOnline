/*
  Script de inicialización de la tienda, correr cuando querramos incializar
*/
IF DB_ID(N'Tienda') IS NULL
    BEGIN
    CREATE DATABASE Tienda
    CONTAINMENT = NONE
    ON PRIMARY
    (
        NAME = Tienda_dat,
        -- Nombre lógico del archivo de datos principal (.mdf)
        FILENAME = 'C:\DATABASES\Tienda.mdf',
        -- Tamaño inicial del archivo de datos
        SIZE = 20MB,
        -- Tamaño máximo permitido
        MAXSIZE = 200MB,
        -- Crecimiento del archivo cuando se llena
        FILEGROWTH = 10MB
    )
    LOG ON
    (
        -- Nombre lógico del archivo de log (.ldf)
        NAME = Calificaciones_log,
        -- Ruta física del log
        FILENAME = 'C:\DATABASES\Tienda.log.ldf',
        -- Tamaño inicial del log
        SIZE = 20MB,
        -- Tamaño máximo del log
        MAXSIZE = 200MB,
        -- Crecimiento del log
        FILEGROWTH = 10MB
    )
    -- Cotejamiento en español moderno, case-insensitive y accent-sensitive
    COLLATE Modern_Spanish_CI_AS;
END;
GO

-- Cambiamos el contexto de ejecución a la base de datos recién creada.
USE Tienda
GO

-- Dropeamos las tablas que ya existen
DROP TABLE IF EXISTS Clientes
DROP TABLE IF EXISTS Carrito
DROP TABLE IF EXISTS Categoria
DROP TABLE IF EXISTS Producto
DROP TABLE IF EXISTS SKU
DROP TABLE IF EXISTS CarritoDetalle
DROP TABLE IF EXISTS Reseña
DROP TABLE IF EXISTS Direccion
DROP TABLE IF EXISTS Pedido
DROP TABLE IF EXISTS Devolucion

CREATE TABLE Clientes
(
    -- Llave primaria, id del usuario
    -- Identity(1,1) es el equivalente a AUTO_INCREMENT
    idUsuario INT NOT NULL IDENTITY(1,1)
    -- Constraint de Primary Key
    CONSTRAINT PK_Cliente PRIMARY KEY,

    -- Campos corrientes
    nombre VARCHAR(50) NOT NULL,
    telefono CHAR(8) NOT NULL,
    correo VARCHAR(30),

    -- Usaremos hashing y salting para crear las contraseñas, esto se hará desde golang
    passwordHash BINARY(32) NOT NULL,
    passwordSalt BINARY(32) NOT NULL,
)

CREATE TABLE Carrito
(
    idCarrito INT IDENTITY(1,1) PRIMARY KEY,
    idUsuario INT UNIQUE,  -- un carrito por cliente
    FOREIGN KEY (idUsuario) REFERENCES Clientes(idUsuario)
    -- Si se elimina el usuario, se elimina el carrito
        ON DELETE CASCADE
);

CREATE TABLE Categoria
(
    idCategoria INT IDENTITY(1,1) PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE Producto
(
    idProducto INT IDENTITY(1,1) PRIMARY KEY,
    descripcion VARCHAR(200) NOT NULL,
    idCategoria INT NULL,

    FOREIGN KEY (idCategoria) REFERENCES Categoria(idCategoria)
    -- Si se borra dejamos la categoria en nulo
        ON DELETE SET NULL
);

CREATE TABLE SKU
(
    idSKU INT IDENTITY(1,1) PRIMARY KEY,
    idProducto INT NOT NULL,
    precio DECIMAL(10,2) NOT NULL CHECK (precio > 0),
    stock INT NOT NULL DEFAULT 0 CHECK (stock > 0),

    FOREIGN KEY (idProducto) REFERENCES Producto(idProducto)
    -- Si se borra el producto, tambien todos los productos minimos vnedibles
        ON DELETE CASCADE
);

CREATE TABLE CarritoDetalle
(
    idDetalle INT IDENTITY(1,1) PRIMARY KEY,
    idCarrito INT NOT NULL,
    idSKU INT NOT NULL,
    cantidad INT NOT NULL CHECK (cantidad > 0),

    FOREIGN KEY (idCarrito) REFERENCES Carrito(idCarrito)
        ON DELETE CASCADE,

    FOREIGN KEY (idSKU) REFERENCES SKU(idSKU)
        ON DELETE CASCADE
);

CREATE TABLE Reseña
(
    idReseña INT IDENTITY(1,1) PRIMARY KEY,
    idUsuario INT NOT NULL,
    idProducto INT NOT NULL,
    puntuación TINYINT CHECK (puntuación BETWEEN 1 AND 5),
    comentario VARCHAR(300),

    FOREIGN KEY (idUsuario) REFERENCES Clientes(idUsuario)
        ON DELETE CASCADE,

    FOREIGN KEY (idProducto) REFERENCES Producto(idProducto)
        ON DELETE CASCADE
);

CREATE TABLE Direccion
(
    idDirección INT IDENTITY(1,1) PRIMARY KEY,
    idUsuario INT NOT NULL,
    tipo VARCHAR(20) CHECK (tipo IN ('Envío','Facturación')),
    detalle VARCHAR(200),

    FOREIGN KEY (idUsuario) REFERENCES Clientes(idUsuario)
        ON DELETE CASCADE
);

CREATE TABLE Pedido
(
    idPedido INT IDENTITY(1,1) PRIMARY KEY,
    idUsuario INT NOT NULL,
    entregado BIT DEFAULT 0,

    FOREIGN KEY(idUsuario) REFERENCES Clientes(idUsuario)
        ON DELETE CASCADE
);

CREATE TABLE Devolucion
(
    idDevolucion INT IDENTITY(1,1) PRIMARY KEY,
    idPedido INT UNIQUE, -- un pedido solo puede tener 1 devolución
    fecha DATE NOT NULL,
    estado VARCHAR(30),
    descripcion VARCHAR(300),
    resolucion VARCHAR(300),

    FOREIGN KEY (idPedido) REFERENCES Pedido(idPedido)
        ON DELETE CASCADE
);
