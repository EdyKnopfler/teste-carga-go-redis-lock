package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

func createRouter(redisSync *redsync.Redsync) *gin.Engine {
	router := gin.Default()

	router.GET("/dotask/:key", func(c *gin.Context) {
		key := c.Param("key")
		mutex := redisSync.NewMutex(key)

		if err := mutex.Lock(); err != nil {
			fmt.Println(err)
			c.JSON(429, gin.H{"message": "barrado no baile"})
			return
		}

		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		n := random.Intn(5000)
		time.Sleep(time.Duration(n) * time.Millisecond)

		if _, err := mutex.Unlock(); err != nil {
			fmt.Println(err)
			c.JSON(429, gin.H{"message": "Erro ao liberar trava"})
			return
		}

		c.JSON(200, gin.H{"message": "pong"})
	})

	return router
}

func appStart() (*http.Server, *goredislib.Client) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	pool := goredis.NewPool(client)
	redisSync := redsync.New(pool)
	router := createRouter(redisSync)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Erro ao criar servidor")
			panic(err)
		}
	}()

	return srv, client
}

func main() {
	srv, client := appStart()
	defer client.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt) // os.Interrupt: Ctrl+C
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("Parando...")

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Erro ao encerrar servidor:", err)
	}

	select {
	case <-ctx.Done():
		fmt.Println("Servidor encerrado")
	}
}
