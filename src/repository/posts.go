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

func (repository Posts) BuscarPostsByID(postID uint64) ([]models.Post, error) {
	rows, err := repository.db.Query(`
		SELECT p.*, u.nick FROM
		posts p INNER JOIN usuarios u
		ON u.id = p.autor_id
		WHERE p.id = ? 
	`, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.Titulo,
			&post.Conteudo,
			&post.AutorID,
			&post.Likes,
			&post.DataCriacao,
			&post.AutorNick,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (repository Posts) BuscarPosts(usuarioID uint64) ([]models.Post, error) {
	rows, err := repository.db.Query(
		`select distinct p.*, u.nick from posts p
		inner join usuarios u on u.id = p.author_id
		inner join posts s on p.author_id = s.usuario_id
		where u.id = ? or s.seguidor_id = ?;
		`,
		usuarioID, usuarioID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		if err = rows.Scan(
			&post.ID,
			&post.Titulo,
			&post.Conteudo,
			&post.AutorID,
			&post.Likes,
			&post.DataCriacao,
			&post.AutorNick,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
