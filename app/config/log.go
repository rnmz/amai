package config

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func InitLogger() (*os.File, error) {
	logDir := "./log"
	logPath := filepath.Join(logDir, fmt.Sprintf("app_%s.log", time.Now().Format("20060102150405")))

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("[Log] Error creating log file: %w", err)
	}
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("[Log] Error opening file: %w", err)
	}

	handler := &PlainFileHandler{
		mu:    &sync.Mutex{},
		out:   file,
		level: slog.LevelInfo,
	}

	slog.SetDefault(slog.New(handler))
	return file, nil
}

type PlainFileHandler struct {
	mu    *sync.Mutex
	out   io.Writer
	level slog.Leveler
}

func (h *PlainFileHandler) Enabled(_ context.Context, level slog.Level) bool {
	if h.level == nil {
		return true
	}
	return level >= h.level.Level()
}

func (h *PlainFileHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	timeStr := r.Time.Format("2006-01-02 15:04:05.000")
	fmt.Fprintf(h.out, "%s [%s] %s", timeStr, r.Level.String(), r.Message)

	r.Attrs(func(a slog.Attr) bool {
		a.Value = a.Value.Resolve()
		v := a.Value.Any()

		if err, ok := v.(error); ok {
			fmt.Fprintf(h.out, " %s=%s", a.Key, err.Error())
		} else {
			fmt.Fprintf(h.out, " %s=%v", a.Key, v)
		}

		return true
	})

	fmt.Fprintln(h.out)
	return nil
}

func (h *PlainFileHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *PlainFileHandler) WithGroup(_ string) slog.Handler {
	return h
}
