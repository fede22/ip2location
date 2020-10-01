//	- Obtener 50 IPs de Argentina (IP, ciudad y pais)
//	- Obtener toda a información de una determinada IP (la IP debe ser un parámetro)
//	- Obtener todos los ISP de Francia (nombres)
//	- Obtener cantidad de IPs por país (país debe ser un parámetro)
//	- Obtener los 3 Proxy Type que más aparecen
//	- Realizar UT.
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {c.JSON(200, gin.H{"message": "pong"})})
	//TODO fifty?
	r.GET("/country/:country/ip/fifty", func(c *gin.Context) {c.JSON(200, gin.H{"message": c.FullPath()})})
	r.GET("/ip/:address", func(c *gin.Context) {c.JSON(200, gin.H{"message": c.FullPath()})})
	//TODO only for France?
	r.GET("/country/:country/isp", func(c *gin.Context) {c.JSON(200, gin.H{"message": c.FullPath()})})
	//TODO count?
	r.GET("/country/:country/ip/count", func(c *gin.Context) {c.JSON(200, gin.H{"message": c.FullPath()})})
	//TODO top?
	r.GET("/top/proxy_type", func(c *gin.Context) {c.JSON(200, gin.H{"message": c.FullPath()})})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
