package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const timeLayout = "2006-01-02T15:04:05.000"

func InitLogger(level, dir string) (*os.File, error) {
	lvl := slog.LevelInfo
	if level != "" {
		if err := lvl.UnmarshalText([]byte(level)); err != nil {
			return nil, fmt.Errorf("[Log] bad level %q: %w", level, err)
		}
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("[Log] create dir %q: %w", dir, err)
	}

	file, err := createLogFile(dir)
	if err != nil {
		return nil, err
	}

	h := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: lvl,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 && a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(timeLayout))
			}
			return a
		},
	})
	slog.SetDefault(slog.New(h))

	return file, nil
}

func createLogFile(dir string) (*os.File, error) {
	name := fmt.Sprintf("app_%s_%d.log", time.Now().Format("20060102_150405"), os.Getpid())
	path := filepath.Join(dir, name)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("[Log] open file %q: %w", path, err)
	}
	return f, nil
}
