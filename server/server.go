// Fazer um hello world
package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro ao criar o servidor: ", err)
		return
	}

	defer ln.Close()

	fmt.Println("Servidor rodando na porta 8080")

	//Loop para aceitar conexões
	for {
		connect, err := ln.Accept()

		if err != nil {
			fmt.Println("Erro ao aceitar conexão: ", err)
			continue
		}

		// Lidar com a conexão em uma goroutine
		go handleConnection(connect)
	}
}

func handleConnection(connect net.Conn) {
	defer func() {
		connect.Close()
		fmt.Println("Conexão fechada: ", connect.RemoteAddr())
	}()

	fmt.Println("Nova conexão estabelecida: ", connect.RemoteAddr())
	buffer := make([]byte, 1024)
	n, err := connect.Read(buffer)
	if err != nil {
		fmt.Println("Erro ao receber mensagem: ", err)
		return
	}
	// Buffer para receber a mensagem do cliente
	fmt.Println("Mensagem recebida pelo cliente: ", string(buffer[:n]))

	_, err = connect.Write([]byte("Olá, cliente!"))
	if err != nil {
		fmt.Println("Erro ao enviar mensagem: ", err)
		return
	}

}
