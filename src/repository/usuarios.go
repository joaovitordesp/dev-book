package repository

import (
	"api-bk/src/models"
	"database/sql"
	"fmt"
)

type usuarios struct {
	db *sql.DB
}

// criar um repositorio de usuarios
func NewRepositoryUsuarios(db *sql.DB) *usuarios {
	return &usuarios{db}
}

func (repository usuarios) Criar(usuario models.Usuario) (uint64, error) {
	statement, err := repository.db.Prepare(
		"insert into usuarios (nome, nick, email, senha) values (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer statement.Close()

	result, err := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if err != nil {
		return 0, err
	}

	lastIDInsert, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastIDInsert), nil
}

func (repository usuarios) Buscar(nomeOuNick string) ([]models.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick) // %nomeOuNick%

	rows, err := repository.db.Query(
		"select id, nome, nick, email, criadoEm from usuarios where nome LIKE ? or nick LIKE ?",
		nomeOuNick, nomeOuNick)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []models.Usuario

	for rows.Next() {
		var usuario models.Usuario

		if err = rows.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, usuario)
	}
	return usuarios, nil
}

func (repository usuarios) BuscarPorID(id uint64) (models.Usuario, error) {
	row, err := repository.db.Query(
		"select id, nome, nick, email, criadoEm from usuarios where id =?", id)

	if err != nil {
		return models.Usuario{}, err
	}

	defer row.Close()
	var usuario models.Usuario

	if row.Next() {
		if err = row.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); err != nil {
			return models.Usuario{}, err
		}
	}
	return usuario, nil
}

func (repository usuarios) AtualizarUsuario(usuarioID uint64, usuario models.Usuario) error {
	statement, err := repository.db.Prepare("update usuarios set nome = ?, nick=?, email=? where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuarioID)

	if err != nil {
		return err
	}

	return nil
}

func (repository usuarios) DeletarUsuario(usuarioID uint64) error {
	statement, err := repository.db.Prepare("delete from usuarios where id =?")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(usuarioID)

	if err != nil {
		return err
	}

	return nil
}

func (repository usuarios) BuscarPorEmail(email string) (models.Usuario, error) {
	row, err := repository.db.Query("select id, senha from usuarios where email =?", email)
	if err != nil {
		return models.Usuario{}, err
	}

	defer row.Close()

	var usuario models.Usuario

	if row.Next() {
		if err = row.Scan(&usuario.ID, &usuario.Senha); err != nil {
			return models.Usuario{}, err
		}
	}

	return usuario, nil
}
