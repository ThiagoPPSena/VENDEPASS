package requests

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Route struct {
	From string `json:"from"`
	To   string `json:"to"`
}

var ServerAddress = "localhost"
var ServerPort = "8080"
var HeaderCpf string

// Coloca um timeout de 2 segundos para a conexão
var ConnectionTimeout = 2 * time.Second

func RequestServer(request string) ([]byte, error) {
	//Conectar ao servidor tcp porta 8080
	connect, err := net.DialTimeout("tcp", ServerAddress+":"+ServerPort, ConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("falha na conexão com o servidor")
	}
	//Garantir que a conexão será fechada
	defer connect.Close()

	//Enviar a string de requisição para o servidor
	_, err = connect.Write([]byte(request))

	if err != nil {
		return nil, fmt.Errorf("falha ao enviar a requisição")
	}

	//Le a resposta da requisição do servidor
	buffer := make([]byte, 1024)
	size, err := connect.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler a resposta do servidor")
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
func StringBuy(routes []Route) string {
	// Cria um slice para armazenar as strings formatadas das rotas
	var routeStrings []string
	for _, route := range routes {
		// Formata cada rota como "origem-destino"
		routeStrings = append(routeStrings, route.From+"/"+route.To)
	}

	// Junta todas as rotas com quebra de linha
	routesJoined := strings.Join(routeStrings, "\n")

	// Cria a string da requisição
	request := "BUY\n" + "COUNT=" + strconv.Itoa(len(routes)) + "\n" + routesJoined
	return request
}
