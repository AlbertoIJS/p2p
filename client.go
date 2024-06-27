package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	const address = "localhost:9999"

	socket, err := net.Dial("udp", address)
	if err != nil {
		log.Fatal("Error al crear socket:", err)
	}
	defer socket.Close()

	// Buscar un archivo en la red
	fileName := "code.txt"
	_, err = socket.Write([]byte(fileName))
	if err != nil {
		fmt.Println("Error al enviar nombre del archivo:", err)
		return
	}

	var data []byte
	buffer := make([]byte, 1024)
	for {
		length, err := socket.Read(buffer)
		if err != nil {
			fmt.Println("Error al leer respuesta:", err)
			return
		}

		// Si la respuesta es "Archivo no encontrado", termina el bucle
		if string(buffer[:length]) == "Archivo no encontrado" {
			fmt.Println(string(buffer[:length]))
			return
		}

		// Agrega los datos leídos al slice de bytes
		data = append(data, buffer[:length]...)

		// Si los datos leídos son menos que el tamaño del buffer, termina el bucle
		if length < len(buffer) {
			break
		}
	}

	// Escribe el slice de bytes en un archivo
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Println("Error al escribir el archivo:", err)
		return
	}

	fmt.Println("Archivo descargado:", fileName)
}
