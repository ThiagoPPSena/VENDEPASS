package graphs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Route struct {
	From  string
	To    string
	Seats int
}

var Graph map[string][]Route

func ReadRoutes() {
	// Abre o arquivo JSON
	file, err := os.Open("./files/routes.json")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	// Decodifica o arquivo JSON
	err = json.NewDecoder(file).Decode(&Graph)
	if err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return
	}
}

// Método para encontrar todas as rotas disponíveis dada uma origem e destino
func FindRoutes(graph map[string][]Route, origin string, destination string, visited map[string]bool, path []string, allpaths *[][]string) {
	visited[origin] = true
	path = append(path, origin)

	// Se a origem for igual ao destino, adiciona a rota encontrada
	if origin == destination {
		tempPath := make([]string, len(path)) // Faz uma cópia do caminho
		copy(tempPath, path)
		*allpaths = append(*allpaths, tempPath)
	} else {
		// Verifica vizinhos (rotas possíveis) e faz a busca recursiva
		for _, neighbor := range graph[origin] {
			if neighbor.Seats > 0 && !visited[neighbor.To] {
				FindRoutes(graph, neighbor.To, destination, visited, path, allpaths)
			}
		}
	}

	// Marca como não visitado (permite outras rotas usarem essa cidade novamente)
	visited[origin] = false
}

func SaveSeats() {

	// Abre o arquivo JSON existente com as opções de sobrescrita
	file, err := os.OpenFile("./files/routes.json", os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close() // Garante que o arquivo será fechado no final

	// Converte os dados para JSON e escreve no arquivo
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Indenta o JSON para uma melhor legibilidade
	err = encoder.Encode(Graph)
	if err != nil {
		fmt.Println("Erro ao escrever os dados no arquivo:", err)
		return
	}

}
