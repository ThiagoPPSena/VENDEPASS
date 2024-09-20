package main

import (
	"VENDEPASS/server/graphs"
	"VENDEPASS/server/responses"
)

func main() {

	graphs.ReadRoutes()
	responses.ReceiveRequest() //Conectando ao cliente

}
