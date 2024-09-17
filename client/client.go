package main

import (
	"VENDEPASS/client/requests"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Response struct {
	Status int                `json:"status"`
	Routes [][]requests.Route `json:"routes"`
}

type ResponseBuy struct {
	Status  int    `json:"status"`
	Message string `json:"routes"`
}

var reader = bufio.NewReader(os.Stdin)

func input() string {
	value, _ := reader.ReadString('\n') // Lê até a quebra de linha
	value = strings.TrimSpace(value)    // Remove espaços extras e a quebra de linha
	return value
}

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

func decodeResponse[T any](response []byte) (T, error) {
	var decodedData T

	err := json.Unmarshal(response, &decodedData)
	if err != nil {
		return decodedData, err
	}
	return decodedData, nil
}

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

func ternaryString(condition bool, trueValue string, falseValue string) string {
	if condition {
		return trueValue
	}
	return falseValue
}

func chooseRoute() (string, string) {
	var origin, destination, response string

	for response != "s" {
		clearConsole()
		// Precisamos de alguma forma listar todas as origens
		fmt.Println("Escolha a origem: ")
		origin = input()
		clearConsole()
		// Precisamos de alguma forma listar todos os destinos
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

func identificationMenu() string {
	var cpf string
	var invalidCpf bool
	for len(cpf) != 11 {
		clearConsole()
		fmt.Println("Faça sua indentificação: ")
		fmt.Print(ternaryString(invalidCpf, "CPF inválido\nDigite um CPF válido: ", "Digite seu CPF: "))
		invalidCpf = false
		fmt.Scan(&cpf)
		if len(cpf) != 11 {
			invalidCpf = true
		}
	}
	return cpf
}

func showTicket(response Response) {
	for i, routeSet := range response.Routes {
		fmt.Println("")
		fmt.Printf("Passagem %d: ", i+1)
		for j, route := range routeSet {
			if len(routeSet)-1 == j {
				fmt.Printf(" Trecho %d: %s -> %s\n", j+1, route.From, route.To)
			} else {
				fmt.Printf("  Trecho %d: %s -> %s", j+1, route.From, route.To)
			}
		}
	}
}

func chooseTicket(response Response) []requests.Route {
	var option int

	fmt.Println("Escolha qual passagem deseja comprar")
	fmt.Scanln(&option)
	selectedTicket := response.Routes[option-1]
	clearConsole()
	for i, route := range selectedTicket {
		fmt.Printf("Trecho %d: %s -> %s ", i+1, route.From, route.To)
	}
	fmt.Printf("\nVocê escolheu a passagem %d, com os trechos acima!", option)
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
			// //Busca no servidor os trechos disponíveis entre a origem e o destino
			response, err := availableTickets(origin, destination)
			if err != nil {
				fmt.Println(err)
				return
			}
			//Mostra ao usuário as rotas disponíveis e deixa ele escolher
			showTicket(response)
			selectedTicket := chooseTicket(response)
			//Quando tiver as rotas escolhidas, chama a função de compra de passagem
			responseBuy, errBuy := buyTicket(selectedTicket)
			if errBuy != nil {
				fmt.Println(errBuy)
				return
			}
			if responseBuy.Status == 200 {
				fmt.Println("Compra efetuada com sucesso")
			} else if responseBuy.Status == 204 {
				fmt.Println(responseBuy.Message)
			} else {
				fmt.Println("Ocorreu um erro inesperado ao efetuar a comprar, cheque suas passagens")
			}
			time.Sleep(5 * time.Second)

		case 2:
			clearConsole()
			availableTickets("Feira de Santana", "São Paulo")
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
	requests.ServerAddress = "192.168.227.102"
	requests.ServerPort = "8080"

	cpf := identificationMenu()
	requests.HeaderCpf = cpf

	defaultMenu()
}
