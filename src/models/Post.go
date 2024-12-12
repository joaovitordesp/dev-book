package models

import (
	"errors"
	"strings"
	"time"
)

type Post struct {
	ID          uint64    `json:"id,omitempty"`
	Titulo      string    `json:"titulo,omitempty"`
	Conteudo    string    `json:"conteudo,omitempty"`
	AutorID     uint64    `json:"autor_id,omitempty"`
	AutorNick   string    `json:"autor_nick,omitempty"`
	Likes       uint64    `json:"likes"`
	DataCriacao time.Time `json:"data_criacao,omitempty"`
}

func (post *Post) Preparar() error {
	if err := post.Validar(); err != nil {
		return err
	}

	post.formatar()
	return nil
}

func (post *Post) Validar() error {
	// Validações de dados
	if post.Titulo == "" {
		return errors.New("o título é obrigatório e não pode estar em branco")
	}
	if post.Conteudo == "" {
		return errors.New("o conteúdo é obrigatório e não pode estar em branco")
	}

	return nil
}

func (post *Post) formatar() {
	post.Titulo = strings.TrimSpace(post.Titulo)
	post.Conteudo = strings.TrimSpace(post.Conteudo)
}
