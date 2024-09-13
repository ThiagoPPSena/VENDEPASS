package requests

import (
	"net"
	"strconv"
	"strings"
)

var ServerAddress = ""
var ServerPort = ""

func RequestServer(request string) string {
	//Conectar ao servidor tcp porta 8080
	connect, err := net.Dial("tcp", ServerAddress+":8080")
	if err != nil {
		panic("Erro ao conectar ao servidor: " + err.Error())
	}
	//Garantir que a conexão será fechada
	defer connect.Close()

	//Enviar a string de requisição para o servidor
	_, err = connect.Write([]byte(request))

	if err != nil {
		panic("Erro ao enviar mensagem: " + err.Error())
	}
	//Le a resposta da requisição do servidor
	buffer := make([]byte, 1024)
	size, err := connect.Read(buffer)
	if err != nil {
		panic("Erro ao ler mensagem: " + err.Error())
	}
	return string(buffer[:size])
}

// GET
// FEIRA DE SANTANA
// SAO PAULO
// Função para gerar a string de requisição
func StringGet(origin string, destination string) string {
	request := "GET\n" + strings.ToUpper(origin) + "\n" + strings.ToUpper(destination)
	return request
}

// BUY
// COUNT=2
// FEIRA DE SANTANA/SALVADOR
// SALVADOR/SAO PAULO
// Função para gerar a srting de requisição de compra
func StringBuy(routes []string) string {
	request := "BUY\n" + "COUNT=" + strconv.Itoa(len(routes)) + "\n" + strings.Join(routes, "\n")
	return request
}
