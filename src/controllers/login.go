package controllers

import (
	"api-bk/src/auth"
	"api-bk/src/database"
	"api-bk/src/models"
	"api-bk/src/repository"
	"api-bk/src/response"
	"api-bk/src/security"
	"encoding/json"
	"io"
	"net/http"
)

func LoginUsuario(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	var usuario models.Usuario
	if err := json.Unmarshal(body, &usuario); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Conectar()
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	defer db.Close()

	repository := repository.NewRepositoryUsuarios(db)
	userSaveOnDatabase, err := repository.BuscarPorEmail(usuario.Email)

	if err != nil {
		response.Erro(w, http.StatusNotFound, err)
		return
	}

	if err = security.VerificarSenha(userSaveOnDatabase.Senha, usuario.Senha); err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	token, err := auth.CriarToken(userSaveOnDatabase.ID)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

}
