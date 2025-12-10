-- Registrar devolucion ligada a un pedido
INSERT INTO Devolucion (idPedido, fecha, estado, descripcion, resolucion)
VALUES (@orderId, @date, @status, @description, @resolution);
