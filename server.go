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

type Node struct {
	Files []File
}

var nodes []Node
var mutex = &sync.Mutex{}

// Búsqueda por índices locales
func searchFile(fileName string) *File {
	mutex.Lock()
	defer mutex.Unlock()

	for _, node := range nodes {
		for _, file := range node.Files {
			if file.Name == fileName {
				return &file
			}
		}
	}

	return nil
}

func main() {
	pc, err := net.ListenPacket("udp", ":9999")
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

	var nodeFiles []File
	for _, entry := range fileEntries {
		if !entry.IsDir() {
			data, err := os.ReadFile("./files/" + entry.Name())
			if err != nil {
				log.Fatal(err)
			}
			nodeFiles = append(nodeFiles, File{Name: entry.Name(), Data: data})
		}
	}
	nodes = append(nodes, Node{Files: nodeFiles})

	buffer := make([]byte, 1024)
	fmt.Println("Server is listening...")
	for {
		length, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}

		fileName := string(buffer[:length])
		file := searchFile(fileName)
		if file != nil {
			// Divide el archivo en chunks de 1024 bytes y los envía uno por uno
			for i := 0; i < len(file.Data); i += 1024 {
				end := i + 1024
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
