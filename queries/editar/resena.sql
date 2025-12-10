-- Editar reseña de un producto
UPDATE [Reseña]
SET idUsuario = @userId,
    idProducto = @productId,
    [puntuación] = @rating,
    comentario = @comment
WHERE idReseña = @id;
