-- Registrar reseña con puntuacion de 1 a 5
INSERT INTO [Reseña] (idUsuario, idProducto, [puntuación], comentario)
VALUES (@userId, @productId, @rating, @comment);
