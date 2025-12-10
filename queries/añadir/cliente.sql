-- Salt para el hash, de esta manera, la misma contraseña no genera siempre el mismo hash y es mas dificil
-- de intificar contraseñas comunes
DECLARE @SALT BINARY(32);
SET @SALT = CRYPT_GEN_RANDOM(32);

-- Valores para insertar a la tabla
INSERT INTO Clientes (nombre, telefono, correo, passwordHash, passwordSalt)
VALUES (
    @name,
    @phone,
    @email,
    -- Hash del salt + la contraseña del usuario
    HASHBYTES('SHA2_256', @SALT + CONVERT(VARBINARY(MAX), @password )),
    @SALT
)