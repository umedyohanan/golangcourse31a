package postgres

import (
	"GoNews/pkg/storage"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			author_id,
			title,
			content,
			created_at
		FROM posts
		ORDER BY id;
	`,
	)
	if err != nil {
		return nil, err
	}
	var posts []storage.Post
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, p)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return posts, rows.Err()
}

func (s *Store) AddPost(post storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO posts (author_id, title, content, created_at)
		VALUES ($1, $2, $3, $4);
		`,
		post.AuthorID,
		post.Title,
		post.Content,
		post.CreatedAt,
	)
	return err
}

func (s *Store) UpdatePost(post storage.Post) error {
	sqlStatement := `
		UPDATE posts 
		SET author_id = $1, title = $2, content = $3
		WHERE id = $4;`
	_, err := s.db.Exec(context.Background(), sqlStatement, post.AuthorID, post.Title, post.Content, post.ID)

	return err
}

func (s *Store) DeletePost(post storage.Post) error {
	sqlStatement := `
		DELETE FROM posts WHERE id = $1;`
	_, err := s.db.Exec(context.Background(), sqlStatement, post.ID)
	return err
}
