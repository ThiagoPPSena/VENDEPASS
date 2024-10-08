package main

import (
	"bufio"
	"client/requests"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)
// Estrutura de resposta para a função de busca de passagens disponíveis
type Response struct {
	Status int                `json:"status"`
	Routes [][]requests.Route `json:"routes"`
}
// Estrutura de resposta para a função de compra de passagens
type ResponseBuy struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
// Estrutura de resposta para a função de busca de todas as passagens
type ResponseGetAll struct {
	Status   int              `json:"status"`
	Passages []requests.Route `json:"passages"`
}
// Cria um leitor de entrada para o console
var reader = bufio.NewReader(os.Stdin)
// Função que lê a entrada do usuário
func input() string {
	value, _ := reader.ReadString('\n') // Lê até a quebra de linha
	value = strings.TrimSpace(value)    // Remove espaços extras e a quebra de linha
	return value
}

func waitForEnter() {
	fmt.Println("Pressione Enter para continuar...")
	reader.ReadString('\n') // Aguarda até que o usuário pressione Enter
}
// Função que limpa o console
func clearConsole() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println("Plataforma não suportada")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
// Função que decodifica a resposta do servidor
func decodeResponse[T any](response []byte) (T, error) {
	var decodedData T

	err := json.Unmarshal(response, &decodedData)
	if err != nil {
		return decodedData, err
	}
	return decodedData, nil
}
// Função que compra a passagem
func buyTicket(routes []requests.Route) (ResponseBuy, error) {
	//Cria uma requisição de compra de rotas
	request := requests.StringBuy(routes)
	//Envia a requisição para o servidor
	response, err := requests.RequestServer(request)
	if err != nil {
		return ResponseBuy{}, err
	}
	data, err := decodeResponse[ResponseBuy](response)
	if err != nil {
		return ResponseBuy{}, err
	}
	return data, nil

}
// Função que busca as passagens disponíveis
func availableTickets(origin string, destination string) (Response, error) {
	// Cria a requisição de compra de passagem
	request := requests.StringGet(origin, destination)
	// Envia a requisição para o servidor
	response, err := requests.RequestServer(request)
	if err != nil {
		return Response{}, err
	}
	data, err := decodeResponse[Response](response)

	if err != nil {
		return Response{}, err
	}
	return data, nil
}
// Função que retorna um valor de acordo com a condição
func ternaryString(condition bool, trueValue string, falseValue string) string {
	if condition {
		return trueValue
	}
	return falseValue
}
// Função que permite o usuário escolher a rota
func chooseRoute() (string, string) {
	var origin, destination, response string

	for response != "s" {
		clearConsole()
		fmt.Println("Escolha a origem: ")
		origin = input()
		clearConsole()
		fmt.Println("Escolha o destino: ")
		destination = input()
		clearConsole()
		fmt.Printf("Origem: %s\nDestino: %s\n", origin, destination)
		fmt.Println("Essa rota está correta? (s/n)")
		response = input()
		clearConsole()
	}
	return origin, destination
}
// Função que mostra o menu de identificação
func identificationMenu() string {
	var cpf string
	var invalidCpf bool
	for len(cpf) != 11 {
		clearConsole()
		fmt.Println("Faça sua indentificação: ")
		fmt.Print(ternaryString(invalidCpf, "CPF inválido\nDigite um CPF válido: ", "Digite seu CPF: "))
		invalidCpf = false
		fmt.Scanln(&cpf)
		if len(cpf) != 11 {
			invalidCpf = true
		}
	}
	return cpf
}
// Função que mostra as passagens disponíveis
func showTicket(response Response) {
	for i, routeSet := range response.Routes {
		fmt.Printf("Passagem %2d: ", i+1)
		for j, route := range routeSet {
			if len(routeSet)-1 == j {
				fmt.Printf("%s -> %s\n", route.From, route.To)
			} else {
				fmt.Printf("%s -> %s / ", route.From, route.To)
			}
		}
	}
}
// Função que mostra as passagens compradas
func showTicketPurchased(passages []requests.Route) {
	fmt.Println("Passagens compradas: ")
	for i, route := range passages {
		fmt.Printf("Passagem %2d: %s -> %s\n", i+1, route.From, route.To)
	}
}
// Função que permite o usuário escolher uma passagem
func chooseTicket(response Response) []requests.Route {
	var option int
	incorrectOption := false
	for {	
		clearConsole()
		showTicket(response)
		fmt.Print(ternaryString(incorrectOption, "Passagem inválida, tente novamente: ", "Escolha uma passagem: "))
		fmt.Scanln(&option)
		if option > 0 && option <= len(response.Routes) {
			break
		}
		incorrectOption = true
	}
	selectedTicket := response.Routes[option-1]
	clearConsole()
	for i, route := range selectedTicket {
		if len(selectedTicket)-1 == i {
			fmt.Printf("%s -> %s\n", route.From, route.To)
		} else {
			fmt.Printf("%s -> %s / ", route.From, route.To)
		}
	}
	fmt.Printf("Você escolheu a passagem %d, com os trechos acima!\n", option)
	return selectedTicket
}

// Roda o menu padrão da interface do cliente
func defaultMenu() {
	var option int
	var invalidOption bool
	for {
		clearConsole()
		fmt.Println("1 - Comprar passagem")
		fmt.Println("2 - Minhas passagens")
		fmt.Println("3 - Sair")
		fmt.Print(ternaryString(invalidOption, "Opção inválida, Escolha uma opção válida: ", "Escollha uma opção: "))
		invalidOption = false
		fmt.Scanln(&option)
		switch option {
		case 1:
			//Primeiro faz o usuario escolher a rota desejada
			var origin, destination string
			origin, destination = chooseRoute()
			// Busca no servidor os trechos disponíveis entre a origem e o destino
			fmt.Println("Buscando passagens disponíveis...")
			response, err := availableTickets(origin, destination)
			clearConsole()
			if err != nil {
				fmt.Println(err)
				waitForEnter()
				continue
			}
			//Mostra ao usuário as rotas disponíveis e deixa ele escolher
			if len(response.Routes) == 0 {
				fmt.Println("Não há passagens disponíveis para essa rota")
				waitForEnter()
				continue
			}
			selectedTicket := chooseTicket(response)
			waitForEnter()
			//Quando tiver as rotas escolhidas, chama a função de compra de passagem
			fmt.Println("Comprando passagem...")
			responseBuy, errBuy := buyTicket(selectedTicket)
			clearConsole()
			if errBuy != nil {
				fmt.Println(errBuy)
				waitForEnter()
				continue
			}
			if responseBuy.Status == 200 {
				fmt.Println("Compra efetuada com sucesso!")
			} else if responseBuy.Status == 204 {
				fmt.Println(responseBuy.Message)
			} else {
				fmt.Println("Ocorreu um erro inesperado ao efetuar a comprar, cheque suas passagens")
			}
			waitForEnter()
		case 2:
			clearConsole()
			fmt.Println("Buscando suas passagens...")
			request := requests.StringGetAll()
			response, err := requests.RequestServer(request)
			if err != nil {
				fmt.Println(err)
				waitForEnter()
				continue
			}
			data, err := decodeResponse[ResponseGetAll](response)
			if err != nil {
				fmt.Println(err)
				waitForEnter()
				continue
			}
			clearConsole()
			if data.Passages == nil {
				fmt.Println("Você não possui passagens compradas")
			} else {
				showTicketPurchased(data.Passages)
			}
			waitForEnter()
		case 3:
			clearConsole()
			fmt.Println("Saindo...")
			os.Exit(0)
		default:
			invalidOption = true
		}
	}
}

func main() {
   // Lê o endereço do servidor e a porta a partir de variáveis de ambiente
	 serverAddress := os.Getenv("SERVER_ADDRESS")
	 serverPort := os.Getenv("SERVER_PORT")
	 if serverAddress != "" {
		requests.ServerAddress = serverAddress
	 }
	 if serverPort != "" {
		requests.ServerPort = serverPort
	 }

	cpf := identificationMenu()
	requests.HeaderCpf = cpf

	defaultMenu()
}
