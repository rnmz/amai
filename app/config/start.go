package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func InitEnvVariables() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

func CheckEnvParams() error {
	vars := []string{
		"TRUSTED_PROXY_IPV4",
		"TRUSTED_PROXY_IPV6",
		"BACKEND_PORT",
		"FILE_PATH",
		"DSN",
	}

	for _, v := range vars {
		if os.Getenv(v) == "" {
			return fmt.Errorf("environment variable %s is not set", v)
		}
	}
	return nil
}

func ShowStartMessage() {
	fmt.Printf("Server started at %s. Start time: %s\n", os.Getenv("BACKEND_PORT"), time.Now().String())
}

func StartServer() {
	db := InitDatabase()
	r := GinApp(db)
	Routing(r)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("BACKEND_PORT"),
		Handler: r.Handler(),

		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,

		MaxHeaderBytes: 1 << 20,
	}

	cleanupTicker := time.NewTicker(5 * time.Second)
	defer cleanupTicker.Stop()

	cleanupDone := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-cleanupTicker.C:
				CleanupSessions()
			case <-cleanupDone:
				return
			}
		}
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("ListenAndServe() error", "error msg", err)
		}
		slog.Info("Server started")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown error. Message %s", err.Error())
	}

	cleanupDone <- true
	CloseDatabase(db)
}
