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

func (repository usuarios) Follow(usuarioID, seguidorID uint64) error {
	statement, err := repository.db.Prepare("insert into seguidores (usuario_id, seguidor_id) values(?,?)")
	if err != nil {
		return err
	}

	defer statement.Close()

	if _, err := statement.Exec(usuarioID, seguidorID); err != nil {
		return err
	}

	return nil
}

func (repository usuarios) Unfollow(usuarioID, seguidorID uint64) error {
	statement, err := repository.db.Prepare("delete from seguidores where usuario_id =? and seguidor_id =?")
	if err != nil {
		return err
	}

	defer statement.Close()

	if _, err := statement.Exec(usuarioID, seguidorID); err != nil {
		return err
	}

	return nil
}

func (repository usuarios) SearchFollowers(usuarioID uint64) ([]models.Usuario, error) {
	rows, err := repository.db.Query("select u.id, u.nome, u.nick, u.email, u.criadoEm from usuarios u join seguidores s on u.id = s.seguidor_id where s.usuario_id =?", usuarioID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var usuarios []models.Usuario
	for rows.Next() {
		var usuario models.Usuario
		if err := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Nick, &usuario.Email, &usuario.CriadoEm); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, usuario)
	}
	return usuarios, nil
}

func (repository usuarios) FollowingUsers(usuarioID uint64) ([]models.Usuario, error) {
	rows, err := repository.db.Query(`
		select u.id, u.nick, u.email, u.criandoEm from usuarios u join seguidores s on u.id s.usuario_id where s.seguidor_id = ?`, usuarioID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var usuarios []models.Usuario

	for rows.Next() {
		var usuario models.Usuario
		if err := rows.Scan(&usuario.ID, &usuario.Nick, &usuario.Email, &usuario.CriadoEm); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (repository usuarios) SavedOnDatabase(usuarioID uint64) (string, error) {
	row, err := repository.db.Query("select senha from usuarios where senha_id=?", usuarioID)
	if err != nil {
		return "", err
	}
	defer row.Close()

	var usuario models.Usuario
	if row.Next() {
		if err := row.Scan(&usuario.Senha); err != nil {
			return "", err
		}
	}
	return usuario.Senha, nil
}
