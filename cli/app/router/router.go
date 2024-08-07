package router

import (
	"html/template"
	"net/http"
	"strings"

	healthcheck "github.com/RaMin0/gin-health-check"
	brotli "github.com/anargu/gin-brotli"
	"github.com/gin-gonic/gin"
	"github.com/semihalev/gin-stats"
)

func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

func staticCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/public") {
			c.Header("Cache-Control", "public, max-age=31536000")
		}
		c.Next()
	}
}

func renderHTMLError(c *gin.Context, status int, message string) {
	c.HTML(status, "pages/error", gin.H{
		"status":  status,
		"message": message,
	})
	c.Abort()
}

func renderJSONError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status":  status,
		"message": message,
	})
	c.Abort()
}

func Router() *gin.Engine {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"safeHTML": safeHTML,
	})

	r.Use(brotli.Brotli(brotli.DefaultCompression))
	r.Use(healthcheck.Default())

	r.NoRoute(func(c *gin.Context) {
		// of the request is to the /api path, return a JSON error
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			renderJSONError(c, http.StatusNotFound, "Page not found")
			return
		}

		renderHTMLError(c, http.StatusNotFound, "Page not found")
	})

	// if GIN_MODE is release, we need to compress static assets and set cache headers
	if gin.Mode() == gin.ReleaseMode {
		r.Use(staticCacheMiddleware())
	}

	r.Static("/public", "./public")

	v1 := r.Group("/api/v1")
	{
		v1.Use(stats.RequestStats())

		v1.GET("/stats", func(c *gin.Context) {
			c.JSON(http.StatusOK, stats.Report())
		})

		devStages := v1.Group("/developmental-stages")
		{
			devStages.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Hello, developmental stages!"})
			})
		}

		neurons := v1.Group("/neurons")
		{
			neurons.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Hello, neurons!"})
			})
		}

		contacts := v1.Group("/contacs")
		{
			contacts.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Hello, contacts!"})
			})
		}

		synapses := v1.Group("/synapses")
		{
			synapses.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Hello, synapses!"})
			})
		}

		nerveRings := v1.Group("/nerve-rings")
		{
			nerveRings.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Hello, nerve rings!"})
			})
		}

		cphates := v1.Group("/cphates")
		{
			cphates.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Hello, cphates!"})
			})
		}

	}

	return r
}
