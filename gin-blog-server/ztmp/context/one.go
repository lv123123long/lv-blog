package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/user", func(c *gin.Context) {
		// 使用 Get 方法获取上下文中的值
		value, exists := c.Get("key")
		if exists {
			c.JSON(200, gin.H{
				"value": value,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Value not found",
			})
		}

		// 使用 MustGet 方法获取上下文中的值
		value, exists = c.MustGet("key").(string)
		if exists {
			c.JSON(200, gin.H{
				"value": value,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Value not found",
			})
		}
	})

	r.Run(":8080")
}

/*
Get方法返回一个布尔值，表示值是否存在。如果值不存在，Get方法会返回false和nil。
MustGet方法会panic，如果值不存在。
因此，如果你确定值一定存在，你可以使用MustGet方法，否则你应该使用Get方法。
*/
