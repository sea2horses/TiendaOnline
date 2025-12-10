-- Actualizar el usuario vinculado al carrito
UPDATE Carrito
SET idUsuario = @userId
WHERE idCarrito = @id;
