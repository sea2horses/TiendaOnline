-- Editar detalles del cliente (por razones de implementación, el reseteo de contraseña es aparte)
UPDATE Clientes
SET nombre = @name, telefono = @phone, correo = @email
WHERE idUsuario = @id