package config

import (
	"amai/blog/app/auth"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

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
			return fmt.Errorf("Environment variable %s is not set", v)
		}
	}
	return nil
}

func visualWidth(s string) int {
	width := utf8.RuneCountInString(s)
	for _, r := range s {
		if r == '🚀' || r == '✨' {
			width++
		}
	}
	return width
}

func renderBox(lines []string, padding int) string {
	maxLen := 0
	for _, l := range lines {
		w := visualWidth(l)
		if w > maxLen {
			maxLen = w
		}
	}

	innerWidth := maxLen + (padding * 2)

	var sb strings.Builder
	sb.WriteString("╔")
	sb.WriteString(strings.Repeat("═", innerWidth))
	sb.WriteString("╗\n")

	for _, line := range lines {
		w := visualWidth(line)
		leftSpace := (innerWidth - w) / 2
		rightSpace := innerWidth - w - leftSpace

		sb.WriteString("║")
		sb.WriteString(strings.Repeat(" ", leftSpace))
		sb.WriteString(line)
		sb.WriteString(strings.Repeat(" ", rightSpace))
		sb.WriteString("║\n")
	}

	sb.WriteString("╚")
	sb.WriteString(strings.Repeat("═", innerWidth))
	sb.WriteString("╝")
	return sb.String()
}

func ShowStartMessage() {
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}
	version := "ALPHA V1.0"
	now := time.Now().Format("2006-01-02 15:04:05")

	logoLines := []string{
		"",
		" █████╗ ███╗   ███╗ █████╗ ██╗",
		"██╔══██╗████╗ ████║██╔══██╗██║",
		"███████║██╔████╔██║███████║██║",
		"██╔══██║██║╚██╔╝██║██╔══██║██║",
		"██║  ██║██║ ╚═╝ ██║██║  ██║██║",
		"╚═╝  ╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝",
		"",
		"✨ Express more than words ✨",
		"",
	}

	infoLines := []string{
		"Server started 🚀",
		fmt.Sprintf("Amai version: %s", version),
		fmt.Sprintf("Port: %s", port),
		fmt.Sprintf("Start time: %s", now),
	}

	fmt.Println(renderBox(logoLines, 4))
	fmt.Println()
	fmt.Println(renderBox(infoLines, 4))
}

func StartServer() {
	logger, err := InitLogger()
	if err != nil {
		return
	}

	defer logger.Close()

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
				auth.CleanupSessions()
			case <-cleanupDone:
				return
			}
		}
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("[Gin] ListenAndServe() error", "error msg", err)
		}
		slog.Info("[Gin] Server started")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("[Gin] Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("[Gin] Server Shutdown error.", "error", err)
	}

	cleanupDone <- true
	CloseDatabase(db)
}
