-- Direccion de envio o facturacion para el usuario
INSERT INTO Direccion (idUsuario, tipo, detalle)
VALUES (@userId, @type, @detail);
