package main

import (
	"api/src/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"

	_ "api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title SocialMedia-API
// @description RESTful API developed in Golang, intended to serve as the backend for a social networking application
// @BasePath /v1
func main() {
	config.Load()
	r := router.Generate()

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	))

	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))

	fmt.Println("Escutando na porta", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
