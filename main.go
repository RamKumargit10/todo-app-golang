package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todo-app/routes"
)

func main() {
	r := routes.SetUpRoutes()
	port := os.Getenv("APP_PORT")
	fmt.Println("Starting server on port:", port)
	serverError := http.ListenAndServe(":"+port, r)
	if serverError != nil {
		log.Fatal(serverError)
	}
}
