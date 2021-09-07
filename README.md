Consumo de endpoints para llenar base de datos Dgraph, creacion de api rest con go y chi router
exponiendo los endpoints con la informacion cargada

Para iniciar, debemos crear un contenedor de docker con Dgraph
docker run --rm -it -p 8000:8000 -p 8080:8080 -p 9080:9080 dgraph/standalone:v21.03.1

Luego corremos go run . -> para iniciar la aplicacion corriendo en el puerto 3000
