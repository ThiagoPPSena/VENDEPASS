package main

import (
	"client/requests"
	"encoding/json"
	"fmt"
)

type Response struct {
	Status int                `json:"status"`
	Routes [][]requests.Route `json:"routes"`
}

type ResponseBuy struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ResponseGetAll struct {
	Status   int              `json:"status"`
	Passages []requests.Route `json:"passages"`
}


func decodeResponse[T any](response []byte) (T, error) {
	var decodedData T
	// fmt.Printf("Decodificando resposta: %v\n", string(response))
	err := json.Unmarshal(response, &decodedData)
	if err != nil {
		return decodedData, err
	}
	return decodedData, nil
}

func main() {
	requests.ServerAddress = "172.16.103.11"
	requests.ServerPort = "8080"
	requests.HeaderCpf = "06417600513"
	request := requests.StringGet("RECIFE", "SALVADOR")
	response, err := requests.RequestServer(request)
	if err != nil {
		return
	}
	data, err := decodeResponse[Response](response)
	if err != nil {
		return
	}

	for i := 0; i < 20; i++ {
		if len(data.Routes) > 3 {
			requestBuy := requests.StringBuy(data.Routes[2])
			reponseBuy, errResponseBuy := requests.RequestServer(requestBuy)
			if errResponseBuy != nil {
				fmt.Printf("Erro ao fazer a requisição: %v\n", errResponseBuy)
				return
			}

			dataBuy, errBuy := decodeResponse[ResponseBuy](reponseBuy)
			if errBuy != nil {
				fmt.Printf("Erro ao decodificar a resposta: %v\n", errBuy)
				return
			}
			fmt.Print(dataBuy.Status, " ")
		} else {
			fmt.Println("Nenhuma rota disponível")
		}
	}

}
