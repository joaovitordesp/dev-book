package rotas

import (
	"api-bk/src/controllers"
	"net/http"
)

var routesPosts = []Rota{
	{
		URI:                "/posts",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CreatePost,
		RequerAutenticacao: true,
	},
	{
		URI:                "/posts",
		Metodo:             http.MethodGet,
		Funcao:             controllers.FindPosts,
		RequerAutenticacao: true,
	},
	{
		URI:                "/posts/{postID}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.FindPostsById,
		RequerAutenticacao: true,
	},
	{
		URI:                "/posts/{id}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.UpdatePost,
		RequerAutenticacao: true,
	},
	{
		URI:                "/posts/{id}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletePost,
		RequerAutenticacao: true,
	},
	{
		URI:                "/usuarios/{usuarioID}/posts",
		Metodo:             http.MethodGet,
		Funcao:             controllers.FindPostByUsuario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/posts/likes",
		Metodo:             http.MethodPost,
		Funcao:             controllers.LikesPost,
		RequerAutenticacao: true,
	},
	{
		URI:                "/posts/dislikes",
		Metodo:             http.MethodPost,
		Funcao:             controllers.DislikesPost,
		RequerAutenticacao: true,
	},
}
