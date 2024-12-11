package controllers

import (
	"api-bk/src/auth"
	"api-bk/src/database"
	"api-bk/src/models"
	"api-bk/src/repository"
	"api-bk/src/response"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	request, erro := io.ReadAll(r.Body)
	if erro != nil {
		response.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.Usuario
	if erro = json.Unmarshal(request, &usuario); erro != nil {
		response.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if err := usuario.Preparar("cadastro"); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repositorio := repository.NewRepositoryUsuarios(db)
	usuario.ID, erro = repositorio.Criar(usuario)
	if erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	response.JSON(w, http.StatusCreated, erro)
}

func FindPosts(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))

	db, erro := database.Conectar()
	if erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	repository := repository.NewRepositoryUsuarios(db)

	usuarios, erro := repository.Buscar(nomeOuNick)

	if erro != nil {
		response.Erro(w, http.StatusBadRequest, erro)
		return
	}
	response.JSON(w, http.StatusOK, usuarios)
}

func FindPostsById(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, err := strconv.ParseUint(parametros["usuarioID"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	repository := repository.NewRepositoryUsuarios(db)

	usuario, erro := repository.BuscarPorID(usuarioID)

	if erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	response.JSON(w, http.StatusOK, usuario)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, err := strconv.ParseUint(parametros["usuarioID"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	usuarioIDToken, err := auth.ExtrairUsuarioID(r)
	if err != nil || usuarioIDToken != usuarioID {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	if usuarioID != usuarioIDToken {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	request, erro := io.ReadAll(r.Body)
	if erro != nil {
		response.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.Usuario
	if erro = json.Unmarshal(request, &usuario); erro != nil {
		response.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("edicao"); erro != nil {
		response.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repository.NewRepositoryUsuarios(db)
	if erro = repositorio.AtualizarUsuario(usuarioID, usuario); erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	response.JSON(w, http.StatusNoContent, nil)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, err := strconv.ParseUint(parametros["usuarioID"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	usuarioIDToken, err := auth.ExtrairUsuarioID(r)
	if err != nil || usuarioIDToken != usuarioID {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	if usuarioID != usuarioIDToken {
		response.Erro(w, http.StatusForbidden, err)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repository.NewRepositoryUsuarios(db)
	if erro = repositorio.DeletarUsuario(usuarioID); erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	response.JSON(w, http.StatusNoContent, nil)
}
