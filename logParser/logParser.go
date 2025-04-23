package logParser

import (
	"fmt"
	"strings"
	"time"
)

func Parse(input *string) (string, error) {
	// Разделяем строку по символу '['
	parts := strings.Split(*input, "[")

	var result []string
	for _, part := range parts {
		// Разделяем каждую часть по символу ']'
		subParts := strings.Split(part, "]")
		for _, subPart := range subParts {
			// Добавляем непустые части в результат
			if strings.TrimSpace(subPart) != "" {
				result = append(result, strings.TrimSpace(subPart))
			}
		}
	}

	// Преобразуем первую часть, если она существует
	if len(result) > 0 {
		// Парсим дату и время
		timeLayout := "Mon Jan 02 15:04:05.999999 2006"
		parsedTime, err := time.Parse(timeLayout, result[0])
		if err != nil {
			fmt.Println("Ошибка при парсинге времени:", err)
		}

		// Форматируем в нужный формат
		formattedTime := parsedTime.Format("15:04:05.000 02.01.2006")
		result[0] = formattedTime
	}

	// Объединяем все части в одну строку
	finalOutput := strings.Join(result, " ")

	// Выводим начальный input и преобразованную строку
	fmt.Print("Начальный input: ")
	fmt.Println(*input)
	fmt.Print("Преобразованная строка: ")
	fmt.Println(finalOutput)

	return finalOutput, nil
}
