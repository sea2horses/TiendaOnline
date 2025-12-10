-- Agregar item al carrito del cliente
INSERT INTO CarritoDetalle (idCarrito, idSKU, cantidad)
VALUES (@cartId, @skuId, @quantity);
