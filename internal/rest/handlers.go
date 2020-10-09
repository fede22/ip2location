package rest

import (
	"github.com/fede22/ip2location/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProxyService interface {
	GetProxy(address string) (domain.Proxy, error)
	GetISPs(countryCode string) ([]string, error)
	GetIPs(countryCode string, limit int) ([]domain.IP, error)
	GetIPCount(countryCode string) (int, error)
	GetTopProxyTypes(limit int) ([]string, error)
}

func GetIPs(s ProxyService) func(c *gin.Context) {
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

func GetISPs(s ProxyService) func(c *gin.Context) {
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

func GetIPCount(s ProxyService) func(c *gin.Context) {
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

func GetProxy(s ProxyService) func(c *gin.Context) {
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

func GetTopProxyTypes(s ProxyService) func(c *gin.Context) {
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
