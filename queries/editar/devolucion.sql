-- Actualizar datos de una devolucion
UPDATE Devolucion
SET idPedido = @orderId,
    fecha = @date,
    estado = @status,
    descripcion = @description,
    resolucion = @resolution
WHERE idDevolucion = @id;
