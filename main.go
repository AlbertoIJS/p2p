package main

import (
	"fmt"
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

func downloadFile(file *File) {
	err := os.WriteFile(file.Name, file.Data, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	fmt.Println("File downloaded:", file.Name)
}

func main() {
	// Añade nodos y archivos a la red
	nodes = append(nodes, Node{Files: []File{{Name: "file1", Data: []byte("data1")}}})
	nodes = append(nodes, Node{Files: []File{{Name: "file2", Data: []byte("data2")}}})

	// Buscar un archivo en la red
	fileName := "file0"
	file := searchFile(fileName)
	if file != nil {
		fmt.Printf("Archivo encontrado: %s\n", file.Name)
		//downloadFile(file)
	} else {
		fmt.Println("Archivo no encontrado")
	}
}
