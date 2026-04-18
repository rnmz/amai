package data

import (
	"context"
	"fmt"
	"io"
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

	db.GetContext(ctx, &fileInfo, "SELECT * FROM files WHERE file_id = $1", id.String())

	file := filepath.Clean(filepath.Join(path, fileInfo.FileId.String()+fileInfo.FileExt))

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return "", err
	}

	return file, nil
}

func UploadFile(db *sqlx.DB, file io.Reader, ext string) (string, error) {
	path := os.Getenv("FILE_PATH")
	generatedFileName := uuid.NewString()

	filePath := filepath.Join(path, generatedFileName+ext)
	dst, fileErr := os.Create(filePath)

	if fileErr != nil {
		return "", fileErr
	}
	defer dst.Close()

	_, copyErr := io.Copy(dst, file)
	if copyErr != nil {
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
		os.Remove(filePath)
		return "", dbErr
	}

	return generatedFileName, nil
}

func DeleteFile(db *sqlx.DB, ctx context.Context, id uuid.UUID) error {
	path := os.Getenv("BACKEND_FILE_DIR")
	tx, txErr := db.BeginTxx(ctx, nil)
	var fileInfo FileEntity

	db.Select(&fileInfo, "SELECT * FROM files WHERE file_id = $1", id)
	file := filepath.Clean(filepath.Join(path, id.String()+fileInfo.FileExt))

	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	res, execErr := tx.Exec("DELETE FROM files WHERE file_id = $1", id)
	if execErr != nil {
		return execErr
	}

	rows, affErr := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("No rows affected. Id %s", id.String())
	}
	if affErr != nil {
		return affErr
	}

	osErr := os.Remove(filepath.Join(path, file))
	if osErr != nil {
		return osErr
	}

	return tx.Commit()
}
