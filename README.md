# Changelog

## [v1.0.0] - 2025-06-06

### ✨ Añadido
- Conexión inicial con Firebase Firestore
- CRUD completo de usuarios
- CRUD completo de tareas
- Registro y login con contraseña encriptada
- Middleware de autenticación con JWT
- Generación de token JWT válido por 10 minutos
- Separación del código en carpetas (`handlers`, `routes`, `middleware`, etc.)

### 🔧 Cambios técnicos
- Se usa bcrypt para encriptación de contraseñas
- Se implementa Firestore como backend NoSQL

### 🛠️ Estructura del proyecto
- Diseño modular y mantenible siguiendo buenas prácticas
