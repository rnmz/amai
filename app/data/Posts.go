package data

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Post struct {
	Id      uuid.UUID `db:"id"`
	Title   string    `db:"title"`
	Poster  string    `db:"poster"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
	Body    string    `db:"body"`
}

func GetAllPages(db *sqlx.DB, ctx context.Context) (int, error) {
	var items int

	slog.Debug("[DB] Requesting total pages count")
	err := db.GetContext(ctx, &items, "SELECT COUNT(*) FROM posts")
	if err != nil {
		slog.Error("[DB] Failed to count posts", "error", err)
		return 0, err
	}

	pagesCount := (items + 50 - 1) / 50
	slog.Debug("[DB] Total pages calculated", "count", pagesCount, "totalItems", items)
	return pagesCount, nil
}

func GetAllPosts(db *sqlx.DB, ctx context.Context, page int) ([]Post, error) {
	var posts []Post
	offset := (page - 1) * 50

	slog.Info("[DB] Requesting posts list", "page", page, "offset", offset)
	err := db.SelectContext(ctx, &posts, "SELECT * FROM posts ORDER BY created DESC LIMIT 50 OFFSET $1", offset)
	if err != nil {
		slog.Error("[DB] Failed to fetch posts list", "page", page, "error", err)
		return nil, err
	}
	return posts, nil
}

func GetPostById(db *sqlx.DB, ctx context.Context, id uuid.UUID) (Post, error) {
	var post Post

	slog.Info("[DB] Requesting post by ID", "id", id)
	err := db.GetContext(ctx, &post, "SELECT * FROM posts WHERE id = $1", id.String())
	if err != nil {
		slog.Error("[DB] Failed to fetch post", "id", id, "error", err)
		return Post{}, err
	}
	return post, nil
}

func AddPost(db *sqlx.DB, ctx context.Context, post Post) error {
	slog.Info("[Transaction] Starting transaction to add post")
	tx, txErr := db.BeginTxx(ctx, nil)
	currentTime := time.Now().UTC()

	if txErr != nil {
		slog.Error("[Transaction] Failed to start transaction", "error", txErr)
		return txErr
	}
	defer tx.Rollback()

	generatedId := uuid.NewString()

	_, err := tx.NamedExecContext(ctx, `
		INSERT INTO posts (id, title, poster, created, updated, body) 
		VALUES (:id, :title, :poster, :created, :updated, :body)`,
		map[string]any{
			"id":      generatedId,
			"title":   post.Title,
			"poster":  post.Poster,
			"created": currentTime,
			"updated": currentTime,
			"body":    post.Body,
		},
	)
	if err != nil {
		slog.Error("[DB] Failed to insert new post metadata", "generatedId", generatedId, "error", err)
		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		slog.Error("[Transaction] Failed to commit add post transaction", "generatedId", generatedId, "error", commitErr)
		return commitErr
	}

	slog.Info("[DB] Post added successfully", "id", generatedId)
	return nil
}

func EditPost(db *sqlx.DB, ctx context.Context, post Post) error {
	slog.Info("[Transaction] Starting transaction to edit post", "id", post.Id)
	tx, txErr := db.BeginTxx(ctx, nil)

	if txErr != nil {
		slog.Error("[Transaction] Failed to start transaction", "id", post.Id, "error", txErr)
		return txErr
	}
	defer tx.Rollback()

	_, err := tx.NamedExecContext(ctx, `
		UPDATE posts 
		SET title = :title, poster = :poster, updated = :updated, body = :body 
		WHERE id = :id`,
		map[string]any{
			"id":      post.Id.String(),
			"title":   post.Title,
			"poster":  post.Poster,
			"updated": time.Now().UTC(),
			"body":    post.Body,
		},
	)
	if err != nil {
		slog.Error("[DB] Failed to update post metadata", "id", post.Id, "error", err)
		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		slog.Error("[Transaction] Failed to commit edit post transaction", "id", post.Id, "error", commitErr)
		return commitErr
	}

	slog.Info("[DB] Post updated successfully", "id", post.Id)
	return nil
}

func DeletePost(db *sqlx.DB, ctx context.Context, postId uuid.UUID) error {
	slog.Info("[Transaction] Starting transaction to delete post", "id", postId)
	tx, txErr := db.BeginTxx(ctx, nil)

	if txErr != nil {
		slog.Error("[Transaction] Failed to start transaction", "id", postId, "error", txErr)
		return txErr
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", postId.String())
	if err != nil {
		slog.Error("[DB] Failed to delete post", "id", postId, "error", err)
		return err
	}

	rows, affErr := res.RowsAffected()
	if affErr != nil {
		slog.Error("[DB] Failed to get affected rows count", "id", postId, "error", affErr)
		return affErr
	}
	if rows == 0 {
		slog.Warn("[DB] No rows affected during post deletion", "id", postId)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		slog.Error("[Transaction] Failed to commit delete post transaction", "id", postId, "error", commitErr)
		return commitErr
	}

	slog.Info("[DB] Post deleted successfully", "id", postId)
	return nil
}
