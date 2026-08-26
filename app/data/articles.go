package data

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ArticleEntity struct {
	Id      uuid.UUID `db:"id"`
	Title   string    `db:"title"`
	Poster  uuid.UUID `db:"poster"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
	Body    string    `db:"body"`
}

func GetAllPagesArticles(db *sqlx.DB, ctx context.Context) (int, error) {
	var items int

	slog.Debug("[DB] Requesting total pages count")
	err := db.GetContext(ctx, &items, "SELECT COUNT(*) FROM articles")
	if err != nil {
		slog.Error("[DB] Failed to count articles", "error", err)
		return 0, err
	}

	pagesCount := (items + 20 - 1) / 20
	slog.Debug("[DB] Total pages calculated", "count", pagesCount, "totalItems", items)
	return pagesCount, nil
}

func GetAllArticles(db *sqlx.DB, ctx context.Context, page int) ([]ArticleEntity, error) {
	var articleEntities []ArticleEntity
	offset := (page - 1) * 20

	slog.Info("[DB] Requesting articles list", "page", page, "offset", offset)
	err := db.SelectContext(ctx, &articleEntities, "SELECT * FROM articles ORDER BY created DESC LIMIT 20 OFFSET $1", offset)
	if err != nil {
		slog.Error("[DB] Failed to fetch articles list", "page", page, "error", err)
		return nil, err
	}
	return articleEntities, nil
}

func GetArticleById(db *sqlx.DB, ctx context.Context, id uuid.UUID) (ArticleEntity, error) {
	var articleEntity ArticleEntity

	slog.Info("[DB] Requesting article by ID", "id", id)
	err := db.GetContext(ctx, &articleEntity, "SELECT * FROM articles WHERE id = $1", id.String())
	if err != nil {
		slog.Error("[DB] Failed to fetch article", "id", id, "error", err)
		return ArticleEntity{}, err
	}
	return articleEntity, nil
}

func AddArticle(db *sqlx.DB, ctx context.Context, articleEntity ArticleEntity) error {
	slog.Info("[DB] Starting transaction to add article")
	tx, txErr := db.BeginTxx(ctx, nil)
	currentTime := time.Now().UTC()

	if txErr != nil {
		slog.Error("[DB] Failed to start transaction", "error", txErr)
		return txErr
	}
	defer tx.Rollback()

	generatedId := uuid.NewString()

	_, err := tx.NamedExecContext(ctx, `
		INSERT INTO articles (id, title, poster, created, updated, body) 
		VALUES (:id, :title, :poster, :created, :updated, :body)`,
		map[string]any{
			"id":      generatedId,
			"title":   articleEntity.Title,
			"poster":  articleEntity.Poster,
			"created": currentTime,
			"updated": currentTime,
			"body":    articleEntity.Body,
		},
	)
	if err != nil {
		slog.Error("[DB] Failed to insert new article metadata", "generatedId", generatedId, "error", err)
		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		slog.Error("[DB] Failed to commit add article transaction", "generatedId", generatedId, "error", commitErr)
		return commitErr
	}

	slog.Info("[DB] Article added successfully", "id", generatedId)
	return nil
}

func EditArticle(db *sqlx.DB, ctx context.Context, articleEntity ArticleEntity) error {
	slog.Info("[DB] Starting transaction to edit article", "id", articleEntity.Id)
	tx, txErr := db.BeginTxx(ctx, nil)

	if txErr != nil {
		slog.Error("[DB] Failed to start transaction", "id", articleEntity.Id, "error", txErr)
		return txErr
	}
	defer tx.Rollback()

	_, err := tx.NamedExecContext(ctx, `
		UPDATE articles 
		SET title = :title, poster = :poster, updated = :updated, body = :body 
		WHERE id = :id`,
		map[string]any{
			"id":      articleEntity.Id.String(),
			"title":   articleEntity.Title,
			"poster":  articleEntity.Poster,
			"updated": time.Now().UTC(),
			"body":    articleEntity.Body,
		},
	)
	if err != nil {
		slog.Error("[DB] Failed to update article metadata", "id", articleEntity.Id, "error", err)
		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		slog.Error("[DB] Failed to commit edit article transaction", "id", articleEntity.Id, "error", commitErr)
		return commitErr
	}

	slog.Info("[DB] Article updated successfully", "id", articleEntity.Id)
	return nil
}

func DeleteArticle(db *sqlx.DB, ctx context.Context, id uuid.UUID) error {
	slog.Info("[DB] Starting transaction to delete article", "id", id)
	tx, txErr := db.BeginTxx(ctx, nil)

	if txErr != nil {
		slog.Error("[DB] Failed to start transaction", "id", id, "error", txErr)
		return txErr
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, "DELETE FROM articles WHERE id = $1", id.String())
	if err != nil {
		slog.Error("[DB] Failed to delete article", "id", id, "error", err)
		return err
	}

	rows, affErr := res.RowsAffected()
	if affErr != nil {
		slog.Error("[DB] Failed to get affected rows count", "id", id, "error", affErr)
		return affErr
	}
	if rows == 0 {
		slog.Warn("[DB] No rows affected during article deletion", "id", id)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		slog.Error("[DB] Failed to commit delete article transaction", "id", id, "error", commitErr)
		return commitErr
	}

	slog.Info("[DB] Article deleted successfully", "id", id)
	return nil
}
