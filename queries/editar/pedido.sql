-- Cambiar estado de entrega o usuario del pedido
UPDATE Pedido
SET idUsuario = @userId,
    entregado = @delivered
WHERE idPedido = @id;
