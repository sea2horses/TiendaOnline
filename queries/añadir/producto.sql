-- Alta de producto con categoria opcional
INSERT INTO Producto (descripcion, idCategoria)
VALUES (@description, @categoryId);
