# closest-cuentadni-store

## Cómo usar
Crear en el directorio principal un archivo .env como el siguiente (ingresando luego del = los valores relevantes):
``` 
LOCALIDAD=
ADDRESS=
URL=
```
En el campo `URL`, es necesario ingresar la dirección del navegador que aparece cuando clickeamos en "Conocé los locales adheridos acá" al abrir un beneficio, por ejemplo, una url válida es `https://www.bancoprovincia.com.ar/cuentadni/buscadores/verdulerias`

Luego podemos correr el programa ya sea corriendo `go run main.go` o generando el ejecutable primero con `go build` y luego ejecutando el binario con `./closest-cuentadni-store`
