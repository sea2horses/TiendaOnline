-- Actualizar descripcion o categoria de un producto
UPDATE Producto
SET descripcion = @description,
    idCategoria = @categoryId
WHERE idProducto = @id;
