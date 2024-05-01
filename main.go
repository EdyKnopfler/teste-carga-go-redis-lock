package main

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

func main() {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})

	pool := goredis.NewPool(client)
	redisSync := redsync.New(pool)

	router := gin.Default()

	router.GET("/dotask/:key", func(c *gin.Context) {
		key := c.Param("key")
		mutex := redisSync.NewMutex(key)
		defer mutex.Unlock()

		n := random.Intn(3)
		time.Sleep(time.Duration(n) * time.Second)

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run(":8080")
}
