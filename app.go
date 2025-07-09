package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// rootHandler executa a lógica original para o caminho "/"
func rootHandler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		http.Error(w, "Erro ao obter hostname", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Mensagem original que você queria preservar
	fmt.Fprintf(w, "<h1>This request was processed by host: %s</h1> - V3", hostname)
}

// proxyHandler executa a lógica de proxy para todos os outros caminhos
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Extrai a URL de destino do caminho da requisição.
	targetURL := strings.TrimPrefix(r.URL.Path, "/")

	// Constrói a URL completa para a requisição de destino.
	fullTargetURL := "https://" + targetURL
	if r.URL.RawQuery != "" {
		fullTargetURL += "?" + r.URL.RawQuery
	}

	log.Printf("Encaminhando requisição para: %s", fullTargetURL)

	// Cria uma nova requisição para a URL de destino.
	proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, fullTargetURL, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar a requisição para o destino: %v", err), http.StatusInternalServerError)
		return
	}

	// Copia os cabeçalhos da requisição original para a nova requisição.
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Envia a requisição.
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao fazer a requisição para o destino: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copia a resposta (cabeçalhos, status, corpo) de volta para o cliente original.
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// mainRouter decide qual handler chamar com base no caminho da URL.
func mainRouter(w http.ResponseWriter, r *http.Request) {
	// Se o caminho for exatamente "/", chama o handler original.
	if r.URL.Path == "/" {
		rootHandler(w, r)
		return // Encerra a função aqui.
	}
	// Para todos os outros caminhos, chama o handler de proxy.
	proxyHandler(w, r)
}

func main() {
	fmt.Println("Servidor proxy iniciado na porta :80")
	// Todas as requisições agora passam pelo nosso roteador principal.
	http.HandleFunc("/", mainRouter)
	log.Fatal(http.ListenAndServe(":80", nil))
}