package responses

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

var ClientAddress = "192.168.1.12"
var ClientPort = "8080"

func ReceiveRequest() {
	ln, err := net.Listen("tcp", ClientAddress+":"+ClientPort)
	if err != nil {
		fmt.Println("Erro ao criar o servidor: ", err)
		return
	}

	defer ln.Close()

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
	defer conn.Close()
	fmt.Println("Nova conexão estabelecida: ", conn.RemoteAddr())

	for {
		// Buffer para receber a mensagem do cliente
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Erro ao receber mensagem: ", err)
			return
		}
		request := string(buffer[:n])
		fmt.Println("Mensagem recebida pelo cliente")

		// Processar a requisição recebida
		proccessRequest(request)

		// Enviar uma resposta ao cliente (opcional)
		_, err = conn.Write([]byte("Requisição processada\n"))
		if err != nil {
			fmt.Println("Erro ao enviar resposta: ", err)
			break
		}
	}
}

func get(origin string, destination string) {
	fmt.Printf("Origem: %s\nDestino: %s\n", origin, destination)
}

func buy(count int, routes []string) {
	i := 0
	for i < count {
		fmt.Printf("Trecho %d: %s\n", i+1, routes[i])
		i++
	}

}

func proccessRequest(request string) {
	requestSepareted := strings.Split(request, "\n")
	method := requestSepareted[0]

	if strings.ToUpper(method) == "GET" {
		origin := requestSepareted[1]
		destination := requestSepareted[2]
		get(origin, destination)
	} else if strings.ToUpper(method) == "BUY" {
		count, _ := strconv.Atoi(strings.Split(requestSepareted[1], "=")[1])
		routes := requestSepareted[2:]
		buy(count, routes)
	}
}
