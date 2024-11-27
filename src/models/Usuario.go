package models

import (
	"errors"
	"strings"
	"time"
)

type Usuario struct {
	ID       uint64    `json:"id,omitempty"`
	Nome     string    `json:"nome"`
	Nick     string    `json:"nick"`
	Email    string    `json:"email"`
	Senha    string    `json:"senha,omitempty"`
	CriadoEm time.Time `json:"criadoEm"`
}

func (usuario *Usuario) Preparar(etapa string) error {
	if err := usuario.validar(etapa); err != nil {
		return err
	}
	usuario.formatar()
	return nil
}

func (usuario *Usuario) validar(etapa string) error {
	if usuario.Nome == "" {
		return errors.New("nome é obrigatório")
	}
	if usuario.Nick == "" {
		return errors.New("nick é obrigatório")
	}
	if usuario.Email == "" {
		return errors.New("email é obrigatório")
	}
	if etapa == "cadastro" && usuario.Senha == "" {
		return errors.New("senha é obrigatório")
	}

	return nil
}

func (usuario *Usuario) formatar() error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Nick = strings.TrimSpace(usuario.Nick)
	usuario.Email = strings.TrimSpace(usuario.Email)

	return nil
}
