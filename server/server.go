package main

import (
	"VENDEPASS/server/graphs"
	"VENDEPASS/server/passages"
	"VENDEPASS/server/responses"
)

func main() {

	graphs.ReadRoutes()        // Pegando rotas de arquivos
	passages.GetPassages()     // Pegando passagens dos clientes
	responses.ReceiveRequest() //Conectando ao cliente

}
