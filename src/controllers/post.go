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

	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	usuarioID, err := auth.ExtrairUsuarioID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var posts models.Post
	if err = json.Unmarshal(bodyRequest, &posts); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	posts.AutorID = usuarioID

	if err = posts.Preparar(); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Conectar()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	defer db.Close()

	repositorio := repository.NewRepositoryPosts(db)
	posts.ID, err = repositorio.CreatePost(posts)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, posts)
}

func FindPosts(w http.ResponseWriter, r *http.Request) {
	usuarioID, err := auth.ExtrairUsuarioID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	repository := repository.NewRepositoryPosts(db)

	usuarios, erro := repository.BuscarPosts(usuarioID)

	if erro != nil {
		response.Erro(w, http.StatusBadRequest, erro)
		return
	}
	response.JSON(w, http.StatusOK, usuarios)
}

func FindPostsById(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	postID, err := strconv.ParseUint(parametros["postID"], 10, 64)
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
	repository := repository.NewRepositoryPosts(db)

	post, erro := repository.BuscarPostsByID(postID)

	if erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	response.JSON(w, http.StatusOK, post)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	usuarioID, err := auth.ExtrairUsuarioID(r)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	parametros := mux.Vars(r)
	postID, err := strconv.ParseUint(parametros["postID"], 10, 64)
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

	repositorio := repository.NewRepositoryPosts(db)
	postSaved, erro := repositorio.BuscarPostsByID(postID)
	if erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if postSaved.AutorID != usuarioID {
		response.Erro(w, http.StatusForbidden, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post models.Post
	if err = json.Unmarshal(bodyRequest, &post); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if err = post.Preparar(); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if erro = repositorio.UpdatePost(postID, post); erro != nil {

	}

	response.JSON(w, http.StatusNoContent, nil)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	postID, err := strconv.ParseUint(parametros["postID"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	postIDToken, err := auth.ExtrairUsuarioID(r)
	if err != nil || postIDToken != postID {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	if postID != postIDToken {
		response.Erro(w, http.StatusForbidden, err)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repository.NewRepositoryPosts(db)
	if erro = repositorio.DeletePost(postID); erro != nil {
		response.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	response.JSON(w, http.StatusNoContent, nil)
}
