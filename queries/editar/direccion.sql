-- Actualizar direccion de un usuario
UPDATE Direccion
SET idUsuario = @userId,
    tipo = @type,
    detalle = @detail
WHERE idDirección = @id;
