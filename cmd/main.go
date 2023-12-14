package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/szuecs/gin-glog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flag.Set("stderrthreshold", "INFO")
	engine := gin.New()
	engine.Use(ginglog.Logger(time.Second))
	engine.Use(gin.Recovery())
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", "8080"),
		Handler: engine,
	}
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(200, "hello")
	})
	glog.Infoln("initializing server")
	// Initializing the server in a goroutine so that it won't block graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Fatalf("error starting server: %s", err)
		}
	}()
	glog.Info("hello deck of card")

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		glog.Fatalf("server forced to shutdown: %s", err)
	}
}
