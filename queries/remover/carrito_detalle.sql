-- Eliminar item de carrito por ID
DELETE FROM CarritoDetalle
WHERE idDetalle = @id;
