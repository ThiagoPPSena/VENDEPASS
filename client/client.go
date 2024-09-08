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

func buyTicket(routes []string) {

	//Cria uma requisição de compra de rotas
	request := requests.StringBuy(routes)
	//Envia a requisição para o servidor
	response, err := requests.RequestServer(request)
	if err != nil {
		println("Erro ao fazer a requisição: ", err.Error())
		return
	}
	println(response)
	
}

func availableTickets(origin string, destination string) {
	// Cria a requisição de compra de passagem
	request := requests.StringGet(origin, destination)
	// Envia a requisição para o servidor
	response, err := requests.RequestServer(request)
	if err != nil {
		println("Erro ao fazer a requisição: ", err.Error())
		return
	}
	println(response)
}

func ternaryString(condition bool, trueValue string, falseValue string) string {
	if condition {
		return trueValue
	}
	return falseValue
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
			clearConsole()
			buyTicket([]string{"Feira de Santana", "São Paulo"})
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

	defaultMenu()
}
