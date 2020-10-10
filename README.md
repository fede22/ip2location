# ip2location

Para inicializar la API localmente ejecutar desde la raíz del proyecto:

```bash
docker build . -t ip2location-app -f build/app/Dockerfile
docker build build/db -t ip2location-db
docker network create ip2location-network
docker run --name ip2location-db --network ip2location-network -e MYSQL_ALLOW_EMPTY_PASSWORD=true -d ip2location-db
docker run --name ip2location-app --network ip2location-network -p 8080:8080 -d ip2location-app 
```

Tener en cuenta al inicializar la APP con el último comando que si la DB aun no cargó toda la data desde los archivos CSV la conexión desde la APP a la misma puede fallar. Esperar un minuto luego de inicializar la DB para evitar esto.

Los endpoints son:

```http
/ping
/country/:country_code/ip -> Obtener 50 IPs del país.
/country/:country_code/isp -> Obtener todos los ISP del país.
/country/:country_code/ip_count -> Obtener la cantidad de IPs del país.
/ip/:address -> - Obtener toda la información de un Proxy asociado a una determinada IP
/top_proxy_types -> Obtener los 3 Proxy Type que más aparecen en todos los países
```

Ejemplos

```http
curl -X GET 'localhost:8080/ping'
curl -X GET 'localhost:8080/country/br/ip'
curl -X GET 'localhost:8080/country/ar/isp'
curl -X GET 'localhost:8080/country/fr/ip_count'
curl -X GET 'localhost:8080/ip/2.10.192.168'
curl -X GET 'localhost:8080/top_proxy_types'
```


