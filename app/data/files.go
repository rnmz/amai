package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const pageSize = 20

var ErrNotFound = errors.New("file not found")

type FileEntity struct {
	FileId  uuid.UUID `db:"file_id"`
	FileExt string    `db:"file_ext"`
}

func GetAllPagesFile(db *sqlx.DB, ctx context.Context) (int, error) {
	var total int
	if err := db.GetContext(ctx, &total, "SELECT COUNT(*) FROM files"); err != nil {
		return 0, fmt.Errorf("count files: %w", err)
	}
	return (total + pageSize - 1) / pageSize, nil
}

func GetAllFiles(db *sqlx.DB, ctx context.Context, page int) ([]FileEntity, error) {
	var files []FileEntity
	offset := (page - 1) * pageSize

	err := db.SelectContext(ctx, &files,
		"SELECT file_id, file_ext FROM files ORDER BY file_id DESC LIMIT $1 OFFSET $2",
		pageSize, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list files (page %d): %w", page, err)
	}
	return files, nil
}

func GetFileById(db *sqlx.DB, ctx context.Context, id uuid.UUID) (string, error) {
	var file FileEntity
	err := db.GetContext(ctx, &file,
		"SELECT file_id, file_ext FROM files WHERE file_id = $1", id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("file %s: %w", id, ErrNotFound)
	}
	if err != nil {
		return "", fmt.Errorf("get file %s: %w", id, err)
	}

	path := storagePath(file.FileId, file.FileExt)
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("file %s is in database but not on disk: %w", id, err)
	}
	return path, nil
}

func UploadFile(db *sqlx.DB, file io.Reader, ext string) (string, error) {
	id := uuid.New()
	path := storagePath(id, ext)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("create storage dir: %w", err)
	}
	if err := writeFile(path, file); err != nil {
		return "", err
	}

	_, err := db.Exec("INSERT INTO files (file_id, file_ext) VALUES ($1, $2)", id.String(), ext)
	if err != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("save file metadata: %w", err)
	}

	return id.String(), nil
}

func DeleteFile(db *sqlx.DB, ctx context.Context, id uuid.UUID) error {
	var ext string
	err := db.GetContext(ctx, &ext,
		"DELETE FROM files WHERE file_id = $1 RETURNING file_ext", id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("file %s: %w", id, ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("delete file %s: %w", id, err)
	}

	if err := os.Remove(storagePath(id, ext)); err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Warn("[Storage] file removed from DB but not from disk", "id", id, "error", err)
	}
	return nil
}

func storagePath(id uuid.UUID, ext string) string {
	return filepath.Join(os.Getenv("FILE_PATH"), id.String()+ext)
}

func writeFile(path string, src io.Reader) error {
	dst, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		_ = dst.Close()
		_ = os.Remove(path)
		return fmt.Errorf("write file: %w", err)
	}

	if err := dst.Close(); err != nil {
		_ = os.Remove(path)
		return fmt.Errorf("close file: %w", err)
	}
	return nil
}
