package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func requestFile(socket net.Conn, fileName string) ([]byte, error) {
	var data []byte
	buffer := make([]byte, 1024)

	_, err := socket.Write([]byte(fileName))
	if err != nil {
		return nil, fmt.Errorf("Error al enviar nombre del archivo: %w", err)
	}

	for {
		length, err := socket.Read(buffer)
		if err != nil {
			return nil, fmt.Errorf("Error al leer respuesta: %w", err)
		}

		// Si la respuesta es "Archivo no encontrado", termina el bucle
		if string(buffer[:length]) == "Archivo no encontrado" {
			return nil, nil
		}

		// Agrega los datos leídos al slice de bytes
		data = append(data, buffer[:length]...)

		// Si los datos leídos son menos que el tamaño del buffer, termina el bucle
		if length < len(buffer) {
			break
		}
	}

	return data, nil
}

func main() {
	// Lista de direcciones de servidores
	serverAddresses := []string{"localhost:9999", "localhost:9998", "localhost:9997"}

	// Buscar un archivo en la red
	fileName := "firma.png"

	// Solicita el archivo a cada servidor
	for _, address := range serverAddresses {
		socket, err := net.Dial("udp", address)
		if err != nil {
			log.Fatal("Error al crear socket:", err)
		}
		defer socket.Close()

		data, err := requestFile(socket, fileName)
		if err != nil {
			fmt.Println(err)
			return
		}

		if data != nil {
			// Escribe el slice de bytes en un archivo
			err = os.WriteFile(fileName, data, 0644)
			if err != nil {
				fmt.Println("Error al escribir el archivo:", err)
				return
			}

			fmt.Println("Archivo descargado:", fileName)
			return
		}
	}

	fmt.Println("Archivo no encontrado en ninguno de los servidores")
}
