package rest

import (
	"github.com/fede22/ip2location/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(s domain.Service) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", Ping)
	r.GET("/country/:country_code/ip", GetIPs(s))
	r.GET("/country/:country_code/isp", GetISPs(s))
	r.GET("/country/:country_code/ip_count", GetIPCount(s))
	r.GET("/ip/:address", GetProxy(s))
	r.GET("/top_proxy_types", GetTopProxyTypes(s))
	return r
}

func Ping(c *gin.Context) {
	c.String(200, "pong")
}

func GetIPs(s domain.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		countryCode, limit := c.Param("country_code"), 50
		ips, err := s.GetIPs(countryCode, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, ips)
	}
}

func GetISPs(s domain.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		countryCode := c.Param("country_code")
		ispNames, err := s.GetISPs(countryCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, ispNames)
	}
}

func GetIPCount(s domain.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		countryCode := c.Param("country_code")
		count, err := s.GetIPCount(countryCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ip_count": count})
	}
}

func GetProxy(s domain.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		address := c.Param("address")
		proxy, err := s.GetProxy(address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, proxy)
	}
}

func GetTopProxyTypes(s domain.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		limit := 3
		proxyTypes, err := s.GetTopProxyTypes(limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, proxyTypes)
	}
}
