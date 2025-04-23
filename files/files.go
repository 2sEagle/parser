package files

import (
	"log"
	"os"
)

func ReadFile(name string) ([]byte, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}
	return data, nil
}
