package main

import (
    "client/requests"
    "encoding/json"
    "fmt"
    "sync"
    "time"
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
// Função principal de teste
func main() {
    // Configuração da requisição
    requests.ServerAddress = "localhost"
    requests.ServerPort = "8080"
    requests.HeaderCpf = "06417600513"
    request := requests.StringGet("RECIFE", "SALVADOR")
    // Faz a requisição ao servidor
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
    // Exibe as rotas disponíveis
    var wg sync.WaitGroup
    // Número de goroutines para simular compras concorrentes
    numGoroutines := 20 
    // Mutex para proteger a variável maxDuration
    var mu sync.Mutex
    // Variável para armazenar o maior tempo de resposta
    var maxDuration time.Duration

    for i := 0; i < numGoroutines; i++ {
        // Adiciona uma goroutine ao grupo de espera
        wg.Add(1)
        go func() {
            // Adiciona uma chamada ao grupo de espera
            defer wg.Done()
            // Faz 50 compras
            for j := 0; j < 50; j++ { // Número de compras por goroutine
                if len(data.Routes) > 9 {
                    start := time.Now()
                    requestBuy := requests.StringBuy(data.Routes[9])
                    responseBuy, errResponseBuy := requests.RequestServer(requestBuy)
                    duration := time.Since(start)

                    mu.Lock()
                    if duration > maxDuration {
                        maxDuration = duration
                    }
                    mu.Unlock()

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
    fmt.Printf("Todas as compras foram processadas. Maior tempo de resposta: %v\n", maxDuration)
}