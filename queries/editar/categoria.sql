-- Cambiar nombre de la categoria
UPDATE Categoria
SET nombre = @name
WHERE idCategoria = @id;
