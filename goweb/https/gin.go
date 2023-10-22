package main

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure" //需要用到这个第三方库
)

func main() {
	secureMiddleware := secure.New(secure.Options{
		//把http://localhost:5678重定向到https://localhost:5678。这个选项其实可以不写，它是默认行为
		SSLRedirect: true,
		SSLHost:     "localhost:5678",
	})
	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)
			if err != nil {
				c.Abort()
				return
			}
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()

	engine := gin.Default()
	engine.Use(secureFunc)

	engine.GET("/", func(c *gin.Context) {
		c.String(200, "欢迎来到HTTP Secure的世界")
	})
	//启动https（http+tls）服务
	engine.RunTLS("localhost:5678", "config/keys/cert.pem", "config/keys/key.pem")
}
