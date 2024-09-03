package main

import (
	"VENDEPASS/client/requests"
)

func main() {
	// Configurando o endereço e a porta do servidor
	requests.ServerAddress = "172.16.103.223"
	requests.ServerPort = "1010"

	// Criar uma requisição de disponibilidade de rotas
	origin := "Feira de Santana"
	destination := "São Paulo"
	request := requests.StringGet(origin, destination)
	response := requests.RequestServer(request)
	println(response)

	// Criar uma requisição de compra de rotas
	routes := []string{"Feira de Santana/Salvador", "Salvador/São Paulo"}
	request = requests.StringBuy(routes)
	response = requests.RequestServer(request)
	println(response)
}
