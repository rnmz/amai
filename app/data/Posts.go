package data

import (
	"context"
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

	err := db.GetContext(ctx, &items, "SELECT COUNT(*) FROM posts")
	if err != nil {
		return 0, err
	}

	pagesCount := (items + 50 - 1) / 50
	return pagesCount, nil
}

func GetAllPosts(db *sqlx.DB, ctx context.Context, page int) ([]Post, error) {
	var posts []Post
	offset := (page - 1) * 50

	err := db.SelectContext(ctx, &posts, "SELECT * FROM posts ORDER BY created DESC LIMIT 50 OFFSET $1", offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPostById(db *sqlx.DB, ctx context.Context, id uuid.UUID) (Post, error) {
	var post Post

	err := db.GetContext(ctx, "SELECT * FROM posts WHERE id = $1", id.String())
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func AddPost(db *sqlx.DB, ctx context.Context, post Post) error {
	tx, txErr := db.BeginTxx(ctx, nil)
	currentTime := time.Now().UTC()

	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	_, err := tx.NamedExec(`
		INSERT INTO posts (id, title, poster, created, updated, body) 
		VALUES (:id, :title, :poster, :created, :updated, :body)`,
		map[string]any{
			"id":      uuid.NewString(),
			"title":   post.Title,
			"poster":  post.Poster,
			"created": currentTime,
			"updated": currentTime,
			"body":    post.Body,
		},
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func EditPost(db *sqlx.DB, ctx context.Context, post Post) error {
	tx, txErr := db.BeginTxx(ctx, nil)

	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	_, err := tx.NamedExec(`
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
		return err
	}

	return tx.Commit()
}

func DeletePost(db *sqlx.DB, ctx context.Context, postId uuid.UUID) error {
	tx, txErr := db.BeginTxx(ctx, nil)

	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	_, err := tx.Exec("DELETE FROM posts WHERE id = $1", postId.String())
	if err != nil {
		return err
	}

	return tx.Commit()
}
