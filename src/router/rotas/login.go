package rotas

import (
	"api-bk/src/controllers"
	"net/http"
)

var rotaLogin = Rota{
	URI:                "/login",
	Metodo:             http.MethodPost,
	Funcao:             controllers.LoginUsuario,
	RequerAutenticacao: false,
}
