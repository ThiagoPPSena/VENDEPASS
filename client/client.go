package main

import (
	"VENDEPASS/client/requests"
	"fmt"
)

func main() {
	// Configurando o endereço e a porta do servidor
	requests.ServerAddress = "localhost"
	requests.ServerPort = "8080"

	response := requests.RequestServer("GETALL\nHEADER 06417600513")
	fmt.Println(response)
}
