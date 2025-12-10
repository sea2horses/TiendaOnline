-- Nuevo pedido del usuario
INSERT INTO Pedido (idUsuario, entregado)
VALUES (@userId, @delivered);
