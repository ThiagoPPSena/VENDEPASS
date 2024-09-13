package main

import (
	"VENDEPASS/client/requests"
	"fmt"
)

func main() {
	// Configurando o endereço e a porta do servidor
	requests.ServerAddress = "localhost"
	requests.ServerPort = "8080"

	response := requests.RequestServer("BUY\nCOUNT=2\nRECIFE/SALVADOR\nSALVADOR/SAO PAULO")
	fmt.Println(response)
}
