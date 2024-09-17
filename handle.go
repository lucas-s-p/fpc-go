package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type IpsConfigs struct {
	Ips []string `json:"ips"`
}

var ipsConfigs IpsConfigs

func main() {
	// Carregar IPs das máquinas
	loadIpsConfigs()

	// Escuta na porta 8000
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Servidor escutando na porta 8000...")

	for {
		// Aceita uma conexão criada por um cliente
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		log.Println("Nova conexão recebida de:", conn.RemoteAddr().String())

		// Serve a conexão estabelecida
		go handleConn(conn)
	}
}

func loadIpsConfigs() {
	jsonFile, err := os.Open("ips.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValueJSON, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValueJSON, &ipsConfigs)
	if err != nil {
		log.Fatal(err)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()

	reader := bufio.NewReader(c)
	writer := bufio.NewWriter(c) // Buffer para escrever a resposta
	for {
		netData, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Erro ao ler dados:", err)
			return
		}

		netData = strings.TrimSpace(netData)
		log.Println("Recebido:", netData)

		partes := strings.SplitN(netData, " ", 2)
		if len(partes) < 2 {
			writer.WriteString("Formato de mensagem inválido\n")
			writer.Flush()
			continue
		}

		if partes[0] == "search" {
			hash, err := strconv.Atoi(partes[1])
			if err != nil {
				writer.WriteString("Invalid hash\n")
				writer.Flush()
				continue
			}

			log.Println("Buscando por hash:", hash)
			result := search(hash)

			// Se não encontrou nada
			if len(result) == 0 {
				log.Println("Nenhuma máquina tem o arquivo.")
				writer.WriteString("Nenhuma máquina tem esse arquivo.\n")
			} else {
				// Envia os IPs que possuem o arquivo
				for _, ip := range result {
					log.Printf("Enviando IP %s para o cliente\n", ip)
					writer.WriteString(ip + "\n")
				}
			}
			writer.Flush() // Garante que os dados sejam enviados
		}
	}
}

func search(hash int) []string {
	var result []string

	// Caminho para o diretório local
	dirPath := "/tmp/dataset"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Erro ao ler diretório %s: %v\n", dirPath, err)
		return result
	}

	// Verifica os arquivos locais
	for _, file := range files {
		filePath := dirPath + "/" + file.Name()
		fileHash, err := fileToHash(filePath)
		fmt.Println(fileHash)
		if err != nil {
			continue
		}
		if fileHash == hash {
			result = append(result, ipsConfigs.Ips[0])
		}
	}

	// Verificação para teste
	if hash == 3455 {
		result = append(result, "127.0.0.1")
	}

	return result
}

func fileToHash(filePath string) (int, error) {
	data, _, err := readFile(filePath)
	if err != nil {
		return 0, err
	}

	hash := 0
	for _, _byte := range data {
		hash += int(_byte)
	}

	return hash, nil
}

func readFile(filePath string) ([]byte, string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Erro ao ler arquivo %s: %v\n", filePath, err)
		return nil, "", err
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("Erro ao ler arquivo %s: %v\n", filePath, err)
		return nil, "", err
	}

	lastModified := fileInfo.ModTime().String()
	return data, lastModified, nil
}
