package main

import (
	"client/requests"
	"encoding/json"
	"fmt"
	"sync"
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
    requests.ServerAddress = "localhost"
    requests.ServerPort = "8080"
    requests.HeaderCpf = "06417600513"
    request := requests.StringGet("RECIFE", "SALVADOR")
    response, err := requests.RequestServer(request)
    if err != nil {
        fmt.Printf("Erro ao fazer a requisição: %v\n", err)
        return
    }
    data, err := decodeResponse[Response](response)
    if err != nil {
        fmt.Printf("Erro ao decodificar a resposta: %v\n", err)
        return
    }

    var wg sync.WaitGroup
    numGoroutines := 40 // Número de goroutines para simular compras concorrentes

    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 500; j++ { // Número de compras por goroutine
                if len(data.Routes) > 9 {
                    requestBuy := requests.StringBuy(data.Routes[9])
                    responseBuy, errResponseBuy := requests.RequestServer(requestBuy)

                    if errResponseBuy != nil {
                        fmt.Printf("Erro ao fazer a requisição: %v\n", errResponseBuy)
                        return
                    }

                    dataBuy, errBuy := decodeResponse[ResponseBuy](responseBuy)
                    if errBuy != nil {
                        fmt.Printf("Erro ao decodificar a resposta: %v\n", errBuy)
                        return
                    }
                    fmt.Print(dataBuy.Status, " ")
                } else {
                    fmt.Println("Nenhuma rota disponível")
                }
            }
        }()
    }

    wg.Wait() // Espera todas as goroutines terminarem
    fmt.Println("Todas as compras foram processadas")
}