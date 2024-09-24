package main

import (
	"server/graphs"
	"server/passages"
	"server/responses"
)

func main() {
	graphs.ReadRoutes()        // Pegando rotas de arquivos
	passages.GetPassages()     // Pegando passagens dos clientes
	responses.ReceiveRequest() //Conectando ao cliente
}
