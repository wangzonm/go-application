package main

import "github.com/gin-gonic/gin"

func main() {
	engine := gin.Default()
	engine.Use(httpTraceWraper())
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

func httpTraceWraper() gin.HandlerFunc {
	return func(c *gin.Context) {
		//ot := opentracing.GlobalTracer()
		// TODO
	}
}
