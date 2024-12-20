package rotas

import (
	"api-bk/src/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type Rota struct {
	URI                string
	Metodo             string
	Funcao             func(http.ResponseWriter, *http.Request)
	RequerAutenticacao bool
}

func Configurar(r *mux.Router) *mux.Router {
	rotas := rotasUsuarios
	rotas = append(rotas, rotaLogin)
	rotas = append(rotas, routesPosts...)

	for _, rota := range rotas {
		if rota.RequerAutenticacao {
			r.HandleFunc(rota.URI, middleware.Logger(middleware.Autenticar(rota.Funcao))).Methods(rota.Metodo)
		} else {
			r.HandleFunc(rota.URI, middleware.Logger(rota.Funcao)).Methods(rota.Metodo)
		}

		r.HandleFunc(rota.URI, rota.Funcao).Methods(rota.Metodo)
	}

	return r
}
