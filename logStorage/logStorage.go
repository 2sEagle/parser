package logStorage

import (
	"fmt"
	"log"
	"os"
	"parser/files"
	"parser/logParser"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"go.ytsaurus.tech/yt/go/schema"
)

type LogString struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"msg"`
}

type SendLogs struct {
	LogStrings []LogString `json:"logStrings"`
	Hostname   string      `json:"hostname"`
	UpdatedAt  string      `json:"updated"`
}

func CreateScheme(logs *SendLogs) (schema.Schema, error) {
	tableSchema, err := schema.Infer(SendLogs{})
	if err != nil {
		return schema.Schema{}, err
	}
	fmt.Println("Inferred struct schema:")
	spew.Fdump(os.Stdout, tableSchema)
	return tableSchema, nil
}

func newLogString(timestamp, level, message string) (*LogString, error) {
	newLogString := &LogString{
		Timestamp: timestamp,
		Level:     level,
		Message:   message,
		// Hostname:  hostname,
		// UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	return newLogString, nil
}

func CreateLogStorage() *[]LogString {
	var logEntries []LogString
	var lastEntry *LogString
	//var hostname, _ = os.Hostname()

	file, err := files.ReadFile("logs.log")
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}

	// Разделяем содержимое на строки
	lines := strings.Split(string(file), "\n")

	for _, line := range lines {
		// Убираем символы перевода строки
		line = strings.ReplaceAll(line, "\r", "")
		line = strings.ReplaceAll(line, "\t", " ")

		if strings.HasPrefix(line, "[") {
			line, _ = logParser.Parse(&line)
		}

		// Разделяем строку на части
		parts := strings.Fields(line)
		if len(parts) < 3 {
			if lastEntry != nil {
				// Если строка не содержит дату и время, добавляем к предыдущему сообщению
				lastEntry.Message += line
			} else {
				log.Printf("Неверный формат строки: %s", line)
			}
			continue
		}

		// Извлекаем время, дату и сообщение
		timePart := parts[0]
		datePart := parts[1]
		levelPart := parts[2]

		if levelPart[len(levelPart)-1] == ':' {
			newLevelPart := levelPart[:len(levelPart)-1]
			levelPart = newLevelPart
		}
		messagePart := strings.Join(parts[3:], " ")

		// // Проверяем, является ли дата сегодняшней
		// if !isToday(datePart) {
		// 	log.Printf("Пропускаем строку с датой: %s", datePart)
		// 	continue
		// }

		// Проверяем, является ли timePart действительным временем
		if !isValidTime(timePart) {
			if lastEntry != nil {
				// Если время недействительно, добавляем к предыдущему сообщению
				lastEntry.Message += " " + line
			} else {
				log.Printf("Неверный формат времени: %s", timePart)
			}
			continue
		}

		// Создаем экземпляр LogEntry
		logEntry, _ := newLogString(timePart+" "+datePart, levelPart, messagePart)
		logEntries = append(logEntries, *logEntry)
		lastEntry = &logEntries[len(logEntries)-1] // Обновляем указатель на последнюю запись
	}

	return &logEntries
}

func isValidTime(timeString string) bool {
	_, err := time.Parse(time.TimeOnly, timeString) // Формат времени: ЧЧ:ММ:СС.миллисекунды
	return err == nil
}

func RetrieveTodaysLogs(logs *[]LogString) []LogString {
	var todaysLogs []LogString
	today := time.Now().Format("02.01.2006") // Формат даты: ДД.ММ.ГГГГ

	for _, logEntry := range *logs {
		if strings.Contains(logEntry.Timestamp, today) {
			todaysLogs = append(todaysLogs, logEntry)
		}
	}
	fmt.Println("Вывод todaysLogs")
	for _, log := range todaysLogs {
		fmt.Println(log)
	}
	fmt.Println()
	return todaysLogs
}

func CreateSendLogs(logs *[]LogString) SendLogs {
	todaysLogs := RetrieveTodaysLogs(logs)
	return SendLogs{
		LogStrings: todaysLogs,
		Hostname:   "localhost",
		UpdatedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

}
