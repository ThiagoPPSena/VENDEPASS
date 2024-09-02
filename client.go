package main

import (
	"fmt"
	"net"
)

func main() {
	//Conectar ao servidor tcp porta 8080
	connect, err := net.Dial("tcp", "172.16.103.223:1010")
	if err != nil {
		panic("Erro ao conectar ao servidor: " + err.Error())
	}
	//Garantir que a conexão será fechada
	defer connect.Close()

	//Ler a mensagem do servidor
	buffer := make([]byte, 1024)
	n, err := connect.Read(buffer)
	if err != nil {
		panic("Erro ao ler mensagem: " + err.Error())
	}
	fmt.Println("Mensagem recebida do servidor: ", string(buffer[:n]))

	//Envia uma mensagem para o servidor
	_, err = connect.Write([]byte("COMPRAR / Feira de Santana -> Salvador"))

	if err != nil {
		fmt.Println("Erro ao enviar mensagem: ", err)
		return
	}
}

// Função para solitar as rotas disponiveis entre duas rotas de origem e destino
func serviceGet(origin string, destination string ) {
	
}

//Função para comprar cada trecho escolhido
func serviceBuy(routes []string) {
	
}


// BUY
// COUNT=2
// FEIRA DE SANTANA/SALVADOR
// SALVADOR/SAO PAULO


// GET
// FEIRA DE SANTANA
// SAO PAULO

