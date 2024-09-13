package responses

import (
	"VENDEPASS/server/graphs"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Definição da estrutura para representar um trecho/passagem
type Route struct {
	From string `json:"From"`
	To   string `json:"To"`
}

// Definição da estrutura para representar a resposta JSON
type Response struct {
	Status int       `json:"status"`
	Routes [][]Route `json:"routes"`
}

func ReceiveRequest() {
	ln, err := net.Listen("tcp", ":8080") // Ouvindo a porta 8080

	if err != nil {
		fmt.Println("Erro ao criar o servidor: ", err)
		return
	}

	defer ln.Close() // Fechar o listen após executar tudo

	fmt.Println("Servidor rodando")

	//Loop para aceitar conexões
	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("Erro ao aceitar conexão: ", err)
			continue
		}

		// Lidar com a conexão em uma goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close() // Fechar conexão após tudo acabar
	fmt.Println("Nova conexão estabelecida: ", conn.RemoteAddr())

	// Buffer para receber a mensagem do cliente
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Erro na leitura do servidor")
	}
	request := string(buffer[:n]) //Requisição do cliente
	fmt.Println("Mensagem recebida pelo cliente")

	// Processar a requisição recebida
	response, err := proccessRequest(request)

	if err != nil {
		fmt.Println("Erro JSON: ", err)
	}

	// Enviar uma resposta ao cliente
	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Erro ao enviar resposta: ", err)
	}
}

// Função que formata a resposta do GET
func formatGetResponse(routes [][]string) ([]byte, error) {
	var formattedRoutes [][]Route
	for _, route := range routes {
		var steps []Route
		for i := 0; i < len(route)-1; i++ {
			steps = append(steps, Route{
				From: route[i],
				To:   route[i+1],
			})
		}
		formattedRoutes = append(formattedRoutes, steps)
	}

	// Criar a resposta com status e as rotas
	response := Response{
		Status: 200,
		Routes: formattedRoutes,
	}

	// Converte a estrutura para JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar JSON: %w", err)
	}

	return jsonResponse, nil
}

// Método GET que recebe a origem e destino do Passageiro e retorna todas as rotas possíveis
func get(origin string, destination string) ([]byte, error) {

	visited := make(map[string]bool) // Lista para mapear se um nó do grafo já foi visitado
	var path []string                // Salva uma rota
	var allPaths [][]string          // Salva todas as rotas possíveis

	// Método para saber todas as rotas possíveis
	graphs.FindRoutes(graphs.Graph, origin, destination, visited, path, &allPaths)

	response, err := formatGetResponse(allPaths) // Formatando a resposta pra envio pro cliente
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar JSON: %w", err)
	}
	// Exibe o JSON resultante
	return response, nil
}

func buy(count int, routes []string) ([]byte, error) {
	// Criar um delay de 5 segundos
	// time.Sleep(5 * time.Second)

	var purchaseList []string
	i := 0
	for i < count {
		origin := strings.Split(routes[i], "/")[0]
		destination := strings.Split(routes[i], "/")[1]

		if graphs.Graph[origin] == nil {
			return []byte("Nenhuma cidade com essa origem"), nil
		}

		for index, route := range graphs.Graph[origin] {
			if route.To == destination {
				fmt.Printf("Destino encontrado: %s\n", destination)
				fmt.Printf("Assentos: %d\n", route.Seats)
				if route.Seats > 0 {
					graphs.Graph[origin][index].Seats -= 1 // Comprou a passagem
					response := fmt.Sprintf("Compra da passagem de %s a %s efetuada\n", origin, destination)
					purchaseList = append(purchaseList, response)
				} else {
					response := fmt.Sprintf("Passagens esgotadas de %s a %s\n", origin, destination)
					purchaseList = append(purchaseList, response)
				}
			}
		}
		i++
	}

	response, err := json.Marshal(purchaseList)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar JSON: %w", err)
	}

	return response, nil

}

func proccessRequest(request string) ([]byte, error) {
	requestSepareted := strings.Split(request, "\n")
	method := requestSepareted[0]

	if strings.ToUpper(method) == "GET" {
		origin := requestSepareted[1]
		destination := requestSepareted[2]
		return get(origin, destination)
	} else if strings.ToUpper(method) == "BUY" {
		count, _ := strconv.Atoi(strings.Split(requestSepareted[1], "=")[1])
		routes := requestSepareted[2:]
		return buy(count, routes)
	}

	return nil, nil
}
