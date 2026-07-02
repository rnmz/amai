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

type FileEntity struct {
	FileId  uuid.UUID `db:"file_id"`
	FileExt string    `db:"file_ext"`
}

func GetFileById(db *sqlx.DB, ctx context.Context, id uuid.UUID) (string, error) {
	path := os.Getenv("FILE_PATH")
	var fileInfo FileEntity

	slog.Info("[DB] Requesting GetFileById", "id", id)

	if err := db.GetContext(ctx, &fileInfo, "SELECT * FROM files WHERE file_id = $1", id.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("[DB] File metadata not found", "id", id)
			return "", err
		}

		slog.Error("[DB] Failed to fetch file metadata", "id", id, "error", err)
		return "", err
	}

	file := filepath.Clean(filepath.Join(path, fileInfo.FileId.String()+fileInfo.FileExt))
	if _, err := os.Stat(file); os.IsNotExist(err) {
		slog.Error("[DB] GetFileById failed", "error", err)
		return "", err
	}

	slog.Debug("[DB] GetFileById completed successfully", "id", id)
	return file, nil
}

func UploadFile(db *sqlx.DB, file io.Reader, ext string) (string, error) {
	path := os.Getenv("FILE_PATH")
	generatedFileName := uuid.NewString()

	slog.Info("[Storage] Uploading new file", "fileId", generatedFileName)

	filePath := filepath.Join(path, generatedFileName+ext)
	dst, fileErr := os.Create(filePath)

	if fileErr != nil {
		slog.Error("[Storage] Failed to create file on disk", "error", fileErr)
		return "", fileErr
	}
	defer dst.Close()

	_, copyErr := io.Copy(dst, file)
	if copyErr != nil {
		slog.Error("[Storage] Failed to write file data", "error", copyErr)
		os.Remove(filePath)
		return "", copyErr
	}

	_, dbErr := db.NamedExec(
		`INSERT INTO files (file_id, file_ext) VALUES (:FileId, :FileExt)`,
		map[string]string{
			"FileId":  generatedFileName,
			"FileExt": ext,
		},
	)
	if dbErr != nil {
		slog.Error("[DB] Failed to save file metadata to database", "error", fileErr)
		os.Remove(filePath)
		return "", dbErr
	}

	return generatedFileName, nil
}
func DeleteFile(db *sqlx.DB, ctx context.Context, id uuid.UUID) error {
	path := os.Getenv("BACKEND_FILE_DIR")

	tx, txErr := db.BeginTxx(ctx, nil)
	if txErr != nil {
		slog.Error("[Transaction] Failed to start transaction", "error", txErr)
		return txErr
	}
	defer tx.Rollback()

	var fileInfo FileEntity

	slog.Info("[Storage] Deleting file", "id", id)

	db.GetContext(ctx, &fileInfo, "SELECT * FROM files WHERE file_id = $1", id.String())

	slog.Debug("[DB] Fetched file metadata", "id", id, "metadata", fileInfo)

	file := filepath.Clean(filepath.Join(path, id.String()+fileInfo.FileExt))

	res, execErr := tx.ExecContext(ctx, "DELETE FROM files WHERE file_id = $1", id.String())
	if execErr != nil {
		slog.Error("[DB] Failed to delete file metadata", "error", execErr)
		return execErr
	}

	rows, affErr := res.RowsAffected()
	if affErr != nil {
		slog.Error("[DB] Failed to get affected rows", "error", affErr)
		return affErr
	}
	if rows == 0 {
		slog.Warn("[DB] No rows affected during deletion", "id", id)
		return fmt.Errorf("no rows affected for id %s", id.String())
	}

	osErr := os.Remove(file)
	if osErr != nil {
		slog.Error("[Storage] Failed to remove file from disk", "error", osErr)
		return osErr
	}

	if commitErr := tx.Commit(); commitErr != nil {
		slog.Error("[Transaction] Failed to commit transaction", "error", commitErr)
		return commitErr
	}

	slog.Info("[Storage] File deleted successfully", "id", id)
	return nil
}
