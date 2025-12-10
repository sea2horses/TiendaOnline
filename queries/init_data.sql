/*
  Datos de ejemplo para la base Tienda.
  Pensado para correr sobre una BD vacía recién creada con queries/init.sql.
*/
USE Tienda;

-- Clientes de ejemplo (hashes/salt dummy de 32 bytes)
DECLARE @salt1 BINARY(32) = 0x0101010101010101010101010101010101010101010101010101010101010101;
DECLARE @salt2 BINARY(32) = 0x0202020202020202020202020202020202020202020202020202020202020202;
DECLARE @salt3 BINARY(32) = 0x0303030303030303030303030303030303030303030303030303030303030303;

DECLARE @hash1 BINARY(32) = 0xA1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1A1;
DECLARE @hash2 BINARY(32) = 0xB2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2B2;
DECLARE @hash3 BINARY(32) = 0xC3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3C3;

DECLARE @cliente1 INT, @cliente2 INT, @cliente3 INT;

INSERT INTO Clientes (nombre, telefono, correo, passwordHash, passwordSalt)
VALUES ('Ana Torres', '55512345', 'ana@example.com', @hash1, @salt1);
SET @cliente1 = SCOPE_IDENTITY();

INSERT INTO Clientes (nombre, telefono, correo, passwordHash, passwordSalt)
VALUES ('Bruno Diaz', '55567890', 'bruno@example.com', @hash2, @salt2);
SET @cliente2 = SCOPE_IDENTITY();

INSERT INTO Clientes (nombre, telefono, correo, passwordHash, passwordSalt)
VALUES ('Carla Ruiz', '55599999', 'carla@example.com', @hash3, @salt3);
SET @cliente3 = SCOPE_IDENTITY();

-- Carritos (uno por cliente)
DECLARE @carrito1 INT, @carrito2 INT;
INSERT INTO Carrito (idUsuario) VALUES (@cliente1);
SET @carrito1 = SCOPE_IDENTITY();
INSERT INTO Carrito (idUsuario) VALUES (@cliente2);
SET @carrito2 = SCOPE_IDENTITY();

-- Categorías
DECLARE @catRopa INT, @catElectro INT;
INSERT INTO Categoria (nombre) VALUES ('Ropa');
SET @catRopa = SCOPE_IDENTITY();
INSERT INTO Categoria (nombre) VALUES ('Electrónica');
SET @catElectro = SCOPE_IDENTITY();

-- Productos
DECLARE @prodCamisa INT, @prodLaptop INT;
INSERT INTO Producto (descripcion, idCategoria) VALUES ('Camisa de algodón', @catRopa);
SET @prodCamisa = SCOPE_IDENTITY();
INSERT INTO Producto (descripcion, idCategoria) VALUES ('Laptop 14"', @catElectro);
SET @prodLaptop = SCOPE_IDENTITY();

-- SKUs
DECLARE @skuCamisaS INT, @skuCamisaM INT, @skuLaptop INT;
INSERT INTO SKU (idProducto, precio, stock) VALUES (@prodCamisa, 19.99, 20);
SET @skuCamisaS = SCOPE_IDENTITY();
INSERT INTO SKU (idProducto, precio, stock) VALUES (@prodCamisa, 21.99, 15);
SET @skuCamisaM = SCOPE_IDENTITY();
INSERT INTO SKU (idProducto, precio, stock) VALUES (@prodLaptop, 799.00, 5);
SET @skuLaptop = SCOPE_IDENTITY();

-- Detalles de carrito
INSERT INTO CarritoDetalle (idCarrito, idSKU, cantidad)
VALUES
    (@carrito1, @skuCamisaS, 2),
    (@carrito1, @skuLaptop, 1),
    (@carrito2, @skuCamisaM, 1);

-- Pedidos
DECLARE @pedido1 INT, @pedido2 INT;
INSERT INTO Pedido (idUsuario, entregado) VALUES (@cliente1, 0);
SET @pedido1 = SCOPE_IDENTITY();
INSERT INTO Pedido (idUsuario, entregado) VALUES (@cliente2, 1);
SET @pedido2 = SCOPE_IDENTITY();

-- Devoluciones (solo sobre el pedido entregado)
INSERT INTO Devolucion (idPedido, fecha, estado, descripcion, resolucion)
VALUES (@pedido2, '2024-01-15', 'Pendiente', 'Teclado defectuoso', NULL);

-- Reseñas
INSERT INTO [Reseña] (idUsuario, idProducto, [puntuación], comentario)
VALUES
    (@cliente1, @prodCamisa, 5, 'Muy cómoda'),
    (@cliente2, @prodLaptop, 3, 'Buena pero con ruido'),
    (@cliente3, @prodCamisa, 4, 'Color agradable');

-- Direcciones
INSERT INTO Direccion (idUsuario, tipo, detalle)
VALUES
    (@cliente1, 'Envío', 'Calle 1 #123'),
    (@cliente1, 'Facturación', 'Calle 1 #123'),
    (@cliente2, 'Envío', 'Av. Central 456'),
    (@cliente3, 'Envío', 'Boulevard Norte 789');
