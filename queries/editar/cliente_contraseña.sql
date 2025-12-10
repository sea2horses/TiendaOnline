-- Debido a la generación aleatoria del hash, el script de cambio de contraseña debe ser aparte en el cliente
DECLARE @SALT BINARY(32);
SET @SALT = CRYPT_GEN_RANDOM(32);

UPDATE Clientes
SET passwordHash = HASHBYTES('SHA2_256', @SALT + CONVERT(VARBINARY(MAX), @password )),
    passwordSalt = @SALT
WHERE idUsuario = @id