package main

import (
	"VENDEPASS/client/requests"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

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

func buyTicket(routes []string) string {

	//Cria uma requisição de compra de rotas
	request := requests.StringBuy(routes)
	//Envia a requisição para o servidor
	response, err := requests.RequestServer(request)
	if err != nil {
		println("Erro ao fazer a requisição: ", err.Error())
		return ""
	}
	return response
	
}

func availableTickets(origin string, destination string) string {
	// Cria a requisição de compra de passagem
	request := requests.StringGet(origin, destination)
	// Envia a requisição para o servidor
	response, err := requests.RequestServer(request)
	if err != nil {
		println("Erro ao fazer a requisição: ", err.Error())
		return ""
	}
	return response
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
		//Precisamos de alguma forma listar todas as origens
		fmt.Println("Escolha a origem: ")
		fmt.Scan(&origin)
		clearConsole()
		//Precisamos de alguma forma listar todos os destinos
		fmt.Println("Escolha o destino: ")
		fmt.Scan(&destination)
	
		fmt.Printf("Origem: %s\nDestino: %s", origin, destination)
		fmt.Println("Essa rota está correta? (s/n)")
		fmt.Scan(&response)
	}

	return origin, destination
}

func buyTicketMenu(routes []string) {
	var option int
	for {
		clearConsole()
		fmt.Println("Escolha uma rota: ")
		fmt.Println(routes)
		fmt.Println("0 - Voltar")
		fmt.Scan(&option)
		if option == 0 {
			return
		}
	}
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

// Roda o menu padrão da interface do cliente
func defaultMenu() {
	var option int
	var invalidOption bool
	for {
		clearConsole()
		fmt.Println("1 - Comprar passagem")
		fmt.Println("2 - Minhas passagens")
		fmt.Println("3 - Sair")
		fmt.Print(ternaryString(invalidOption, "Opção inválida, Escolha uma opção válida", "Escollha uma opção: "))
		invalidOption = false
		fmt.Scan(&option)
		switch option {
		case 1:
			// Primeiro faz o usuario escolher a rota desejada
			// var origin, destination string
			// origin, destination = chooseRoute()
			// Busca no servidor os trechos disponíveis entre a origem e o destino
			// Precisa fazer a função que retorna as rotas disponíveis, transforma a resposta mais amigavel
			// response := availableTickets(origin, destination)
			// Mostra ao usuário as rotas disponíveis e deixa ele escolher
			// buyTicketMenu(response)
			// Quando tiver as rotas escolhidas, chama a função de compra de passagem
			// buyTicket(routes)
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
	requests.ServerAddress = "127.0.0.1"
	requests.ServerPort = "8080"
	// Pede para o usuário digitar o cpf
	cpf := identificationMenu()
	requests.HeaderCpf = cpf

	defaultMenu()
}
