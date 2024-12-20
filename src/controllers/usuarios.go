package controllers

import (
	"api-bk/src/auth"
	"api-bk/src/database"
	"api-bk/src/models"
	"api-bk/src/repository"
	"api-bk/src/response"
	"api-bk/src/security"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
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

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
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

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
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

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
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

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
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

func SeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, err := auth.ExtrairUsuarioID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	parametros := mux.Vars(r) //lê-se os parametros
	usuarioID, err := strconv.ParseUint(parametros["usuarioID"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if seguidorID == usuarioID {
		response.Erro(w, http.StatusForbidden, errors.New("você não pode seguir você mesmo"))
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repository := repository.NewRepositoryUsuarios(db)
	if err := repository.Follow(seguidorID, usuarioID); err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func ParaDeSeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, err := auth.ExtrairUsuarioID(r)

	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	parametros := mux.Vars(r) //lê-se os parametros
	usuarioID, err := strconv.ParseUint(parametros["usuarioID"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if seguidorID == usuarioID {
		response.Erro(w, http.StatusForbidden, errors.New("você não pode parar de seguir você mesmo"))
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryUsuarios(db)
	if err := repository.Unfollow(seguidorID, usuarioID); err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func SearchFollowers(w http.ResponseWriter, r *http.Request) {
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
	seguidores, err := repository.SearchFollowers(usuarioID)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, seguidores)
}

func FollowingUsers(w http.ResponseWriter, r *http.Request) {
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
	followUsers, err := repository.FollowingUsers(usuarioID)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, followUsers)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	usuarioIDToken, err := auth.ExtrairUsuarioID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseUint(parametros["usuarioID"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if usuarioIDToken != usuarioID {
		response.Erro(w, http.StatusUnauthorized, errors.New("você não tem permissão para alterar este usuário"))
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)

	var password models.Password

	if err = json.Unmarshal(bodyRequest, &password); err != nil {
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

	passSavedOnDatabase, err := repository.SavedOnDatabase(usuarioID)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.VerificarSenha(passSavedOnDatabase, password.CurrentPassword); err != nil {
		response.Erro(w, http.StatusUnauthorized, errors.New("senha atual inválida"))
		return
	}

	passWithHash, err := security.Hash(password.CurrentPassword)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if err = repository.UpdatePassword(usuarioID, string(passWithHash)); err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
