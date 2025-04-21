package main

import (
	"fmt"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		http.Error(w, "Erro ao obter hostname", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body {
					background-color: green;
					color: white;
					font-family: Arial, sans-serif;
					text-align: center;
					padding-top: 50px;
				}
			</style>
		</head>
		<body>
			<h1>This request was processed by host: %s</h1>
			<p> - V1</p>
		</body>
		</html>
		`, hostname)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":80", nil)
}
