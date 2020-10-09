//	- Obtener 50 IPs de Argentina (IP, ciudad y pais)
//	- Obtener toda a información de una determinada IP (la IP debe ser un parámetro)
//	- Obtener todos los ISP de Francia (nombres)
//	- Obtener cantidad de IPs por país (país debe ser un parámetro)
//	- Obtener los 3 Proxy Type que más aparecen
//	- Realizar UT.
package main

import (
	"github.com/fede22/ip2location/internal/domain"
	"github.com/fede22/ip2location/internal/http/rest"
	"github.com/fede22/ip2location/internal/storage/mysql"
	"log"
)

func main() {
	repo, err := mysql.NewRepository()
	if err != nil {
		log.Fatal(err)
	}
	s := domain.NewService(repo)
	r := rest.NewRouter(s)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
