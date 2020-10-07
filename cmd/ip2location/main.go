//	- Obtener 50 IPs de Argentina (IP, ciudad y pais)
//	- Obtener toda a información de una determinada IP (la IP debe ser un parámetro)
//	- Obtener todos los ISP de Francia (nombres)
//	- Obtener cantidad de IPs por país (país debe ser un parámetro)
//	- Obtener los 3 Proxy Type que más aparecen
//	- Realizar UT.
package main

import (
	"github.com/fede22/ip2location/internal/rest"
	"github.com/fede22/ip2location/internal/proxy"
	"github.com/fede22/ip2location/internal/storage/mysql"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()

	repo, err := mysql.NewProxyRepository()
	if err != nil {
		log.Fatal(err)
	}
	s := proxy.NewService(repo)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.GET("/country/:country_code/ip", rest.GetIPs(s))
	r.GET("/country/:country_code/isp", rest.GetISPs(s))
	r.GET("/country/:country_code/ip_count", rest.GetIPCount(s))
	r.GET("/ip/:address", rest.GetProxy(s))
	r.GET("/top_proxy_types", rest.GetTopProxyTypes(s))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
