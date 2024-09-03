package main

import (
	"VENDEPASS/server/responses"
)

func main() {

	// Configurando o endereço e a porta do servidor
	responses.ClientAddress = "192.168.1.12"
	responses.ClientPort = "8080"

	responses.ReceiveRequest() //Conectando ao cliente

}
