package passages

import (
	"encoding/json"
	"fmt"
	"os"
)
// Estrutura de passagens
type MyPassages struct {
	From string
	To   string
}
// Mapa de passagens
var Passages = map[string][]MyPassages{}
// Função para ler as passagens do arquivo JSON
func GetPassages() {
	// Abre o arquivo JSON em modo de leitura
	file, err := os.Open("./files/myPassages.json")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close() // Garante que o arquivo será fechado no final

	// Decodifica o JSON do arquivo para a variável data
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Passages)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo JSON:", err)
		return
	}
}
// Função para salvar as passagens no arquivo JSON
func SavePassages() {
	// Abre o arquivo JSON existente com as opções de sobrescrita
	file, err := os.OpenFile("./files/myPassages.json", os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close() // Garante que o arquivo será fechado no final

	// Converte os dados para JSON e escreve no arquivo
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Indenta o JSON para uma melhor legibilidade
	err = encoder.Encode(Passages)
	if err != nil {
		fmt.Println("Erro ao escrever os dados no arquivo:", err)
		return
	}

}
