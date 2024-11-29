package main

import (
	"api-bk/src/config"
	"api-bk/src/router"
	"fmt"
	"log"
	"net/http"
)

// func init() { //Gerador de  SECRET KEY, deve-se usar só uma vez para incluir o valor no arquivo .env
// 	key := make([]byte, 64)

// 	if _, err := rand.Read(key); err != nil {
// 		log.Fatal(err)
// 	}

// 	stringBase64 := base64.StdEncoding.EncodeToString(key)
// 	fmt.Println("Key do JWT:", stringBase64)
// }

func main() {
	config.Carregar()

	fmt.Println(config.Porta)
	fmt.Println("Rodando API!")

	r := router.Gerar()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r))

	// golang:golang@/devbook?charset=utf8&parseTime=true&loc=Local
}
