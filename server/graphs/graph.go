package graphs

type Route struct {
	From  string
	To    string
	Seats int
}

var Graph = map[string][]Route{
	"SALVADOR":    {{From: "SALVADOR", To: "SAO PAULO", Seats: 50}, {From: "SALVADOR", To: "RECIFE", Seats: 50}},
	"RECIFE":      {{From: "RECIFE", To: "SAO PAULO", Seats: 50}, {From: "RECIFE", To: "SALVADOR", Seats: 1}, {From: "RECIFE", To: "JOAO PESSOA", Seats: 50}},
	"SAO PAULO":   {{From: "SAO PAULO", To: "SALVADOR", Seats: 50}},
	"JOAO PESSOA": {{From: "JOAO PESSOA", To: "ARACAJU", Seats: 50}, {From: "JOAO PESSOA", To: "SALVADOR", Seats: 50}},
	"ARACAJU":     {{From: "ARACAJU", To: "SALVADOR", Seats: 50}},
} // Grafo de voos

// Método para econtrar todas as rotas disponíveis dada uma origem e destino
func FindRoutes(graph map[string][]Route, origin string, destination string, visited map[string]bool, path []string, allpaths *[][]string) {
	visited[origin] = true
	path = append(path, origin)

	if origin == destination {

		*allpaths = append(*allpaths, path)
		visited[origin] = false
		path = nil
		return
	}

	for _, neighbor := range graph[origin] {
		if neighbor.Seats > 0 && !visited[neighbor.To] {
			FindRoutes(graph, neighbor.To, destination, visited, path, allpaths)
		}
	}

}
