package requests

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var ServerAddress = "localhost"
var ServerPort = "8080"
var HeaderCpf string

// Coloca um timeout de 2 segundos para a conexão
var ConnectionTimeout = 10 * time.Second

func RequestServer(request string) ([]byte, error) {
	//Conectar ao servidor tcp porta 8080
	print("Conectando ao servidor...", ServerAddress+":"+ServerPort)
	connect, err := net.DialTimeout("tcp", ServerAddress+":"+ServerPort, ConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao servidor")
	}
	//Garantir que a conexão será fechada
	defer connect.Close()

	//Enviar a string de requisição para o servidor
	_, err = connect.Write([]byte(request))

	if err != nil {
		return nil, fmt.Errorf("erro ao enviar a requisição")
	}

	//Le a resposta da requisição do servidor
	buffer := make([]byte, 1024)
	size, err := connect.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("erro ao receber a resposta")
	}
	return buffer[:size], nil
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
