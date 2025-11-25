package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fawwns/antibruteforce/internal/config"
	"github.com/fawwns/antibruteforce/internal/httpapi"
	"github.com/fawwns/antibruteforce/internal/logger"
	"github.com/fawwns/antibruteforce/internal/service"
)

func main() {
	logger.Init()
	logger.Info.Println("Starting AntiBruteforce service...")

	cfg := config.Load()
	logger.Info.Printf("Loaded config: %+v", cfg)

	svc := service.New(cfg)
	h := httpapi.NewHandler(svc)
	router := httpapi.NewRouter(h)
	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	logger.Info.Printf("Server starting on :%s", cfg.Port)
	go func() {
		for {
			time.Sleep(time.Minute)
			svc.Cleanup()
		}
	}()

	go func() {
		logger.Info.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error.Printf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	logger.Info.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error.Printf("Error during server shutdown: %v", err)
	}

	logger.Info.Println("Server gracefully stopped")
}
