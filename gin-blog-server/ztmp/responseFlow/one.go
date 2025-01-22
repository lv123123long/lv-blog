package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/user", func(c *gin.Context) {
		user := struct {
			Name    string `json:"name"`
			Age     int    `json:"age"`
			Address struct {
				City    string `json:"city"`
				Country string `json:"country"`
			} `json:"address"`
			Emails []string `json:"emails"`
		}{
			Name: "John Doe",
			Age:  30,
			Address: struct {
				City    string `json:"city"`
				Country string `json:"country"`
			}{
				City:    "New York",
				Country: "USA",
			},
			Emails: []string{"john@example.com", "john.doe@example.com"},
		}

		c.JSON(200, user)
	})

	r.Run(":8080")
}




/*
func(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello, World!"},)
}
c.JSON 有三个参数
status：HTTP状态码
data: 要返回的数据

当用户访问/hello时，路由处理函数会返回一个JSON响应，其中包含一个消息字段。

当我们运行这个应用程序并访问http://localhost:8080/hello时，Gin框架会执行路由处理函数，并将JSON响应写入HTTP响应流。然后，Gin框架会处理HTTP响应，将其发送回客户端。

在这个过程中，我们没有显式地返回响应。相反，我们直接通过c.JSON方法将响应写入HTTP响应流。Gin框架会处理HTTP响应，将其发送回客户端。

这就是Gin的响应流的工作方式。路由处理函数不需要显式地返回响应，而是直接将响应写入HTTP响应流。Gin框架会处理HTTP响应，将其发送回客户端。
*/