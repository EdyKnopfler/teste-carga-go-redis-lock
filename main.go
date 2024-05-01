package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

func main() {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})

	pool := goredis.NewPool(client)

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run(":8080")
}
