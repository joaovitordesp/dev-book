package repository

import (
	"api-bk/src/models"
	"database/sql"
)

type Posts struct {
	db *sql.DB
}

func NewRepositoryPosts(db *sql.DB) *Posts {
	return &Posts{db}
}

func (repository Posts) CreatePost(post models.Post) (uint64, error) {
	statement, err := repository.db.Prepare("insert into posts (titulo, conteudo, autor_id) values (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer statement.Close()

	result, err := statement.Exec(post.Titulo, post.Conteudo, post.AutorID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}
