package config

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartRestServer(listenAddr string, port int) {
	r := gin.Default()

	r.GET("/templates", func(c *gin.Context) {
		templateList := make([]string, 0)
		for name := range templates {
			templateList = append(templateList, name)
		}
		c.JSON(http.StatusOK, templateList)
	})

	r.GET("/configs", func(c *gin.Context) {
		configs, err := globalStorage.GetAll()
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, configs)
	})

	r.GET("/configs/:id", func(c *gin.Context) {
		id := c.Param("id")
		config, err := globalStorage.Get(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(config))
	})

	// Add Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run(fmt.Sprintf("%s:%d", listenAddr, port))
}
