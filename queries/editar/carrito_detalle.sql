-- Ajustar items en el carrito
UPDATE CarritoDetalle
SET idCarrito = @cartId,
    idSKU = @skuId,
    cantidad = @quantity
WHERE idDetalle = @id;
