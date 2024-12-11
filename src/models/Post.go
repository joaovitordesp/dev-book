package models

import "time"

type Post struct {
	ID          uint64    `json:"id,omitempty"`
	Titulo      string    `json:"titulo,omitempty"`
	Conteudo    string    `json:"conteudo,omitempty"`
	AutorID     uint64    `json:"autor_id,omitempty"`
	AutorNick   string    `json:"autor_nick,omitempty"`
	Likes       uint64    `json:"likes"`
	DataCriacao time.Time `json:"data_criacao,omitempty"`
}
