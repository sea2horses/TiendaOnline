-- Ajustar datos del SKU (precio y stock)
UPDATE SKU
SET idProducto = @productId,
    precio = @price,
    stock = @stock
WHERE idSKU = @id;
