# postgresqlgo-utils-test
utiles para los test con testcontainers y go

### ¿cómo funciona?
en el paquete **postgresqlgoutilstest** están los siguientes métodos:

- Setup: permite levantar docker para las pruebas con postgres
- LoadScript: realiza la carga de un script en la base de datos de pruebas
- LoadScriptChan: carga el contenido de un script pero enviando batchs de inserts

para funcionar utiliza el proyecto **github.com/sebastian4j/postgresqlgo** que entrega las conexiónes, **github.com/testcontainers/testcontainers-go** para levantar postgresql con docker, **github.com/jackc/pgx/v5** para las conexiones a postgres y **github.com/stretchr/testify** en los test.

;)