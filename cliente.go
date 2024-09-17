package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Conecta ao servidor
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	// Solicitar ao usuário o hash do arquivo
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Digite o hash do arquivo a ser pesquisado: ")
		hash, _ := reader.ReadString('\n')
		hash = strings.TrimSpace(hash)

		// Envia a pesquisa para o servidor
		fmt.Fprintf(conn, "search "+hash+"\n")

		// Espera e lê a resposta do servidor
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Erro ao receber a resposta do servidor:", err)
			return
		}

		// Mostra a resposta
		fmt.Println("Resposta do servidor:", message)
	}
}
