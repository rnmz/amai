package data

import (
	"context"
	"fmt"
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
	var total int
	if err := db.GetContext(ctx, &total, "SELECT COUNT(*) FROM articles"); err != nil {
		return 0, fmt.Errorf("count articles: %w", err)
	}
	return (total + pageSize - 1) / pageSize, nil
}

func GetAllArticles(db *sqlx.DB, ctx context.Context, page int) ([]ArticleEntity, error) {
	var articles []ArticleEntity
	offset := (page - 1) * pageSize

	err := db.SelectContext(ctx, &articles,
		"SELECT id, title, poster, created, updated FROM articles ORDER BY created DESC, id LIMIT $1 OFFSET $2",
		pageSize, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list articles (page %d): %w", page, err)
	}
	return articles, nil
}

func GetArticleById(db *sqlx.DB, ctx context.Context, id uuid.UUID) (ArticleEntity, error) {
	var article ArticleEntity
	err := db.GetContext(ctx, &article,
		"SELECT * FROM articles WHERE id = $1", id.String())
	if err != nil {
		return ArticleEntity{}, fmt.Errorf("get article %s: %w", id, err)
	}
	return article, nil
}

func AddArticle(db *sqlx.DB, ctx context.Context, article ArticleEntity) error {
	now := time.Now().UTC()

	_, err := db.ExecContext(ctx,
		`INSERT INTO articles (id, title, poster, created, updated, body)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.NewString(), article.Title, article.Poster.String(), now, now, article.Body,
	)
	if err != nil {
		return fmt.Errorf("add article: %w", err)
	}
	return nil
}

func EditArticle(db *sqlx.DB, ctx context.Context, article ArticleEntity) error {
	_, err := db.ExecContext(ctx,
		`UPDATE articles
		 SET title = $1, poster = $2, updated = $3, body = $4
		 WHERE id = $5`,
		article.Title, article.Poster.String(), time.Now().UTC(), article.Body, article.Id.String(),
	)
	if err != nil {
		return fmt.Errorf("edit article %s: %w", article.Id, err)
	}
	return nil
}

func DeleteArticle(db *sqlx.DB, ctx context.Context, id uuid.UUID) error {
	if _, err := db.ExecContext(ctx, "DELETE FROM articles WHERE id = $1", id.String()); err != nil {
		return fmt.Errorf("delete article %s: %w", id, err)
	}
	return nil
}
