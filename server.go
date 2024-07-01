package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

type File struct {
	Name string
	Data []byte
}

type Server struct {
	Files []File
}

var servers []Server
var mutex = &sync.Mutex{}

// Búsqueda por índices locales
func searchFile(fileName string) *File {
	mutex.Lock()
	defer mutex.Unlock()

	for _, server := range servers {
		for _, file := range server.Files {
			if file.Name == fileName {
				return &file
			}
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: server <port>")
		return
	}

	port := os.Args[1]
	pc, err := net.ListenPacket("udp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pc.Close()

	// Añade nodos y archivos a la red
	fileEntries, err := os.ReadDir("./files")
	if err != nil {
		log.Fatal(err)
	}

	// Crea instancias de servidores
	numServers := 3
	for i := 0; i < numServers; i++ {
		servers = append(servers, Server{})
	}

	// Distribuye los archivos entre los servidores
	for i, entry := range fileEntries {
		if !entry.IsDir() {
			data, err := os.ReadFile("./files/" + entry.Name())
			if err != nil {
				log.Fatal(err)
			}
			serverIndex := i % numServers
			servers[serverIndex].Files = append(servers[serverIndex].Files, File{Name: entry.Name(), Data: data})
		}
	}

	buffer := make([]byte, 1024)
	fmt.Println("Server is listening on port", port)
	for {
		length, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}

		fileName := string(buffer[:length])
		file := searchFile(fileName)
		maxPacketSize := 1024
		if file != nil {
			// Divide el archivo en chunks de 65535 bytes y los envía uno por uno
			for i := 0; i < len(file.Data); i += maxPacketSize {
				end := i + maxPacketSize
				if end > len(file.Data) {
					end = len(file.Data)
				}
				_, err = pc.WriteTo(file.Data[i:end], addr)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		} else {
			_, err = pc.WriteTo([]byte("Archivo no encontrado"), addr)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
