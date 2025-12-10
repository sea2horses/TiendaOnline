-- Crear un carrito nuevo para un usuario (restriccion UNIQUE en idUsuario)
INSERT INTO Carrito (idUsuario)
VALUES (@userId);
