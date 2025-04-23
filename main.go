package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"parser/logStorage"
	"time"

	"go.ytsaurus.tech/yt/go/ypath"
	"go.ytsaurus.tech/yt/go/yt"
	"go.ytsaurus.tech/yt/go/yt/ythttp"
)

const (
	cluster string = "jupiter.yt.idzn.ru"
	//tokenFile    string = "~.yt/token.txt"
)

func main() {
	logs := logStorage.CreateLogStorage()
	// for _, log := range *logs {
	// 	fmt.Println(log)
	// }
	sendLogs := logStorage.CreateSendLogs(logs)
	fmt.Println("Вывод SendLog")
	for _, sendLog := range sendLogs.LogStrings {
		fmt.Println(sendLog)
	}
	fmt.Println()

	//sendToYT(logs)
	if err := sendToYT(&sendLogs); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		os.Exit(1)
	}
}

func sendToYT(sendLogs *logStorage.SendLogs) error {
	yc, err := ythttp.NewClient(&yt.Config{
		Proxy:             cluster,
		ReadTokenFromFile: true,
		//Token:             tokenFile,
	})
	if err != nil {
		return err
	}

	ctx := context.Background()
	//tablePath := ypath.Path("//home/hc/foresight/logs/go-table-example-" + guid.New().String())
	tablePath := ypath.Path("//home/hc/foresight/logs/" + time.Now().Format("2006-01-02")) //<append=false>

	tableSchema, err := logStorage.CreateScheme(&logStorage.SendLogs{})
	if err != nil {
		log.Fatalf("Ошибка при создании схемы: %v", err)
	}

	_, err = yt.CreateTable(ctx, yc, tablePath, yt.WithSchema(tableSchema))
	if err != nil {
		fmt.Printf("Table at https:%s/navigation?path=%s already exists\n", cluster, tablePath.String())
	}
	fmt.Printf("Created table at https:%s/navigation?path=%s\n", cluster, tablePath.String())

	// Запись данных в таблицу
	writer, err := yc.WriteTable(ctx, tablePath, nil)
	if err != nil {
		return err
	}

	fmt.Println("Writing rows to table...")
	for _, v := range sendLogs.LogStrings {
		if err = writer.Write(v); err != nil {
			return err
		}
	}
	if err = writer.Commit(); err != nil {
		return err
	}
	fmt.Printf("Written and committed %v rows\n", len(sendLogs.LogStrings))

	return nil
}
