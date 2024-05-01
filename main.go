package main

import (
	"context"
	"fmt"
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
		Addr:     "localhost:6379",
		Password: "$3nh4!",
	})

	pool := goredis.NewPool(client)
	redisSync := redsync.New(pool)

	router := gin.Default()

	router.GET("/dotask/:key", func(c *gin.Context) {
		key := c.Param("key")
		ctx := context.Background()
		mutex := redisSync.NewMutex(key)

		if err := mutex.LockContext(ctx); err != nil {
			fmt.Println(err)
			c.JSON(429, gin.H{"message": "barrado no baile"})
			return
		}

		n := random.Intn(500)
		time.Sleep(time.Duration(n) * time.Millisecond)

		if _, err := mutex.UnlockContext(ctx); err != nil {
			fmt.Println(err)
			c.JSON(429, gin.H{"message": "Erro ao liberar trava"})
			return
		}

		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run(":8080")
}
