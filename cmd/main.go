package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/deck/internal/app/config"
	"github.com/deck/internal/app/handler"
	"github.com/deck/internal/app/repo"
	"github.com/deck/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/szuecs/gin-glog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port     string
	logLevel string
)

func main() {
	loadEnvVars()
	_ = flag.Set("stderrthreshold", logLevel)

	engine := gin.New()
	engine.Use(ginglog.Logger(time.Second))
	engine.Use(gin.Recovery())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: engine,
	}

	db, err := config.NewDbConnection()
	if err != nil {
		glog.Fatalf("couldn't connect to db", err.Error())
		panic(err)
	}

	deckRepo := repo.NewDeckRepo(db)
	deckService := service.NewDeckService(deckRepo)
	deckHandler := handler.NewDeckHandler(deckService)
	deckHandler.InitRoutes(engine)

	listenAndServe(server, db)
}

func listenAndServe(server *http.Server, db *sqlx.DB) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	glog.Infoln("initializing server")
	// Initializing the server in a goroutine so that it won't block graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Fatalf("error starting server: %s", err)
		}
	}()

	<-ctx.Done()
	stop()
	glog.Infof("received signal, closing")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		glog.Fatalf("server forced to shutdown: %s", err)
	}
	if err := db.Close(); err != nil {
		glog.Fatalf("couldn't close db: %s", err)
	}
}

func loadEnvVars() {
	err := godotenv.Load(".env")
	if err != nil {
		glog.Fatal("Error loading .env file")
	}

	port = os.Getenv("PORT")
	logLevel = os.Getenv("LOG_LEVEL")
}
