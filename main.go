package main

import (
	"api-bk/src/config"
	"api-bk/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.Carregar()

	fmt.Println(config.Porta)
	fmt.Println("Rodando API!")

	r := router.Gerar()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r))

	// golang:golang@/devbook?charset=utf8&parseTime=true&loc=Local
}
