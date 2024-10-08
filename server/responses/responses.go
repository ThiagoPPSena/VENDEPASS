package responses

import (
	"encoding/json"
	"fmt"
	"net"
	"server/graphs"
	"server/passages"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Definição da estrutura para representar um trecho/passagem
type Route struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Definição da estrutura para representar a resposta JSON
type ResponseGet struct {
	Status int       `json:"status"`
	Routes [][]Route `json:"routes"`
}
// Definição da estrutura para representar a resposta JSON
type ResponseBuy struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
// Definição da estrutura para representar a resposta JSON
type ResponseGetAll struct {
	Status   int                   `json:"status"`
	Passages []passages.MyPassages `json:"passages"`
}
// Função para receber requisições
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
// Função para lidar com a conexão
func handleConnection(conn net.Conn) {
	defer conn.Close() // Fechar conexão após tudo acabar
	fmt.Println("Nova conexão estabelecida: ", conn.RemoteAddr())

	// Buffer para receber a mensagem do cliente
	buffer := make([]byte, 3072)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Erro na leitura do servidor")
	}
	fmt.Println("Mensagem recebida pelo cliente")
	request := string(buffer[:n]) //Requisição do cliente

	// Processar a requisição recebida
	response, err := proccessRequest(request)

	if err != nil {
		fmt.Println("Erro JSON: ", err)
	}

	// Enviar uma resposta ao cliente
	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Erro ao enviar resposta: ")
	}
}

// Função que formata a resposta do GET ALL
func formatGetAllResponse(passages []passages.MyPassages) ([]byte, error) {
	response := ResponseGetAll{
		Status:   200,
		Passages: passages,
	}

	// Converte a estrutura para JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar JSON: %w", err)
	}

	return jsonResponse, nil
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
	response := ResponseGet{
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

	// Ordenando da menor rota para a maior
	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	// Pegando as 10 menores rotas disponíveis
	if len(allPaths) >= 10 {
		allPaths = allPaths[:10]
	}

	response, err := formatGetResponse(allPaths) // Formatando a resposta pra envio pro cliente
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar JSON: %w", err)
	}
	// Exibe o JSON resultante
	return response, nil
}
// Mutex para proteger a compra de passagens
var mutex = &sync.Mutex{}
// Método para comprar passagens
func buy(count int, routes []string, cpf string) ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()

	purchaseMap := make(map[string]int)
	i := 0
	// Repetições para comprar as passagens
	for i < count {
		origin := strings.Split(routes[i], "/")[0]
		destination := strings.Split(routes[i], "/")[1]
		// Verifica se a cidade de origem existe
		if graphs.Graph[origin] == nil {
			response := ResponseBuy{
				Message: fmt.Sprintf("Nenhuma cidade com essa origem: %s", origin),
				Status:  404,
			}

			jsonResponse, err := json.Marshal(response)
			if err != nil {
				return nil, fmt.Errorf("erro ao gerar JSON: %w", err)
			}

			return jsonResponse, nil
		}

		routeFound := false
		// Verifica se a rota existe
		for index, route := range graphs.Graph[origin] {
			if route.To == destination {
				routeFound = true
				if route.Seats > 0 { // Se tiver assento disponível
					purchaseMap[origin] = index
				} else {
					response := ResponseBuy{
						Message: fmt.Sprintf("Passagens esgotadas de %s a %s", origin, destination),
						Status:  204,
					}

					jsonResponse, err := json.Marshal(response)
					if err != nil {
						return nil, fmt.Errorf("erro ao gerar JSON: %w", err)
					}

					return jsonResponse, nil
				}
			}
		}

		if !routeFound {
			response := ResponseBuy{
				Message: fmt.Sprintf("Nenhuma rota encontrada de %s a %s", origin, destination),
				Status:  404,
			}

			jsonResponse, err := json.Marshal(response)
			if err != nil {
				return nil, fmt.Errorf("erro ao gerar JSON: %w", err)
			}
			return jsonResponse, nil
		}
		i++
	}
	// Comprando as passagens
	for key, value := range purchaseMap {
		// Comprando as passagens
		graphs.Graph[key][value].Seats -= 1

		// Salvando as passagens do cliente
		destination := graphs.Graph[key][value].To
		// Criando um novo objeto MyPassages
		newPassage := passages.MyPassages{
			From: key,
			To:   destination,
		}

		// Adicionando o objeto ao map na chave cpf
		passages.Passages[cpf] = append(passages.Passages[cpf], newPassage)
	}

	// Persiste as passagens compradas em arquivo JSON
	passages.SavePassages()
	// Persiste a compra dos assentos em arquivo JSON
	graphs.SaveSeats()

	response := ResponseBuy{
		Message: "Passagens compradas com sucesso",
		Status:  200,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar JSON: %w", err)
	}

	return jsonResponse, nil

}
// Método para obter todas as passagens de um cliente
func getall(cpf string) ([]byte, error) {

	myPassages := passages.Passages[cpf]

	// Convertendo o map para JSON
	myPassagesFormatted, err := formatGetAllResponse(myPassages)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar JSON: %w", err)
	}
	return myPassagesFormatted, nil
}
// Função para processar a requisição
func proccessRequest(request string) ([]byte, error) {
	requestSepareted := strings.Split(request, "\n")
	method := requestSepareted[0]

	if strings.ToUpper(method) == "GET" {
		origin := requestSepareted[1]
		destination := requestSepareted[2]
		return get(origin, destination)
	} else if strings.ToUpper(method) == "BUY" {
		cpf := strings.Split(requestSepareted[1], " ")[1]
		count, _ := strconv.Atoi(strings.Split(requestSepareted[2], "=")[1])
		routes := requestSepareted[3:]
		return buy(count, routes, cpf)
	} else if strings.ToUpper(method) == "GETALL" {
		cpf := strings.Split(requestSepareted[1], " ")[1]
		return getall(cpf)
	}

	return nil, nil
}
