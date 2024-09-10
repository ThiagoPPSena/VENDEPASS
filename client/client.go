package main

import (
	"VENDEPASS/client/requests"
)

func main() {
	// Configurando o endereço e a porta do servidor
	//Passar ip local
	requests.ServerAddress = "127.0.0.1"
	requests.ServerPort = "8080"

	// Criar uma requisição de disponibilidade de rotas
	origin := "Feira de Santana"
	destination := "São Paulo"
	request := requests.StringGet(origin, destination)

	for i := 0; i < 200; i++ {
		_, err := requests.RequestServer(request)
		if err != nil {
			println("Erro ao fazer a requisição: ", err.Error())
			return
		}
	}



	// println(response)

	// Criar uma requisição de compra de rotas
	// routes := []string{"Feira de Santana/Salvador", "Salvador/São Paulo"}
	// request = requests.StringBuy(routes)
	// response = requests.RequestServer(request)
	// println(response)
}
