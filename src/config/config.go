package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	StringConexaoBanco = ""
	Porta              = 5000
	SecretKey          []byte
)

// vai inicializar as variaveis de ambiente
func Carregar() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal("Não foi possível carregar as variaveis de ambiente.")
	}

	Porta, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		Porta = 9000
		http.ListenAndServe(fmt.Sprintf(":%d", Porta), nil)

	}

	StringConexaoBanco = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USUARIO"),
		os.Getenv("DB_SENHA"),
		os.Getenv("DB_NOME"),
	)

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
