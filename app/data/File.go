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

func GetAllPagesFile(db *sqlx.DB, ctx context.Context) (int, error) {
	var items int

	slog.Debug("[DB] Requesting total pages count")
	err := db.GetContext(ctx, &items, "SELECT COUNT(*) FROM files")
	if err != nil {
		slog.Error("[DB] Failed to count files", "error", err)
		return 0, err
	}

	pagesCount := (items + 20 - 1) / 20
	slog.Debug("[DB] Total pages calculated", "count", pagesCount, "totalItems", items)
	return pagesCount, nil
}

func GetAllFiles(db *sqlx.DB, ctx context.Context, page int) ([]FileEntity, error) {
	var files []FileEntity
	offset := (page - 1) * 20

	slog.Info("[DB] Requesting posts list", "page", page, "offset", offset)
	err := db.SelectContext(ctx, &files, "SELECT * FROM files ORDER BY file_id DESC LIMIT 20 OFFSET $1", offset)
	if err != nil {
		slog.Error("[DB] Failed to fetch files list", "page", page, "error", err)
		return nil, err
	}
	return files, nil
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
		`INSERT INTO files (file_id, file_ext) VALUES (:file_id, :file_ext)`,
		FileEntity{
			FileId:  uuid.MustParse(generatedFileName),
			FileExt: ext,
		},
	)
	if dbErr != nil {
		slog.Error("[DB] Failed to save file metadata to database", "error", dbErr)
		os.Remove(filePath)
		return "", dbErr
	}

	return generatedFileName, nil
}

func DeleteFile(db *sqlx.DB, ctx context.Context, id uuid.UUID) error {
	path := os.Getenv("FILE_PATH")

	tx, txErr := db.BeginTxx(ctx, nil)
	if txErr != nil {
		slog.Error("[Transaction] Failed to start transaction", "error", txErr)
		return txErr
	}
	defer tx.Rollback()

	var fileInfo FileEntity
	if err := tx.GetContext(ctx, &fileInfo, "SELECT * FROM files WHERE file_id = $1", id.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("[DB] File metadata not found for deletion", "id", id)
			return fmt.Errorf("file not found: %s", id.String())
		}
		slog.Error("[DB] Failed to fetch file metadata", "id", id, "error", err)
		return err
	}

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

	if commitErr := tx.Commit(); commitErr != nil {
		slog.Error("[Transaction] Failed to commit transaction", "error", commitErr)
		return commitErr
	}

	file := filepath.Clean(filepath.Join(path, id.String()+fileInfo.FileExt))
	if osErr := os.Remove(file); osErr != nil {
		slog.Error("[Storage] Failed to remove file from disk after DB commit", "id", id, "error", osErr)
		return osErr
	}

	slog.Info("[Storage] File deleted successfully", "id", id)
	return nil
}
