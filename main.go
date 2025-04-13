package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"telegram/ai"
	"telegram/config"
	"telegram/tg"
)

func commandHandler(logInfo *log.Logger, logErr *log.Logger, ch chan tg.Message) {
	var command string

	for {
		fmt.Println("Ready to process your command(q - quit, s - start dialog request, c - continue dialog reauest, r - reapeat last message, u - custom dialog request)")
		logInfo.Println("Ready to process messages or command")
		fmt.Scan(&command)
		switch command {
		case "r":
			logInfo.Println("The command to repeat the last message has been received")
			ch <- tg.Message{"Повтори своё предыдущее сообщение", -1}
		case "q":
			logInfo.Panicln(`Введена команда "q"`)
		case "s":
			logInfo.Println("The command to start the dialogue has been received")
			ch <- tg.Message{"Представь, что этого сообщения нет, просто начни диалог", -1}
		case "c":
			logInfo.Println("The command to continue the dialogue has been received")
			ch <- tg.Message{"Представь, что этого сообщения нет, просто продолжи диалог", -1}
		case "u":
			logInfo.Println("The command to custom request has been received")
			var customMsg string
			fmt.Println("Введите запрос для сообщения")
			in := bufio.NewReader(os.Stdin)
			customMsg, _ = in.ReadString('\n')
			logInfo.Println(`Custom request contained "` + customMsg + `"`)
			ch <- tg.Message{customMsg, -1}
		}
	}
}

func dialog(login, password string, logInfo *log.Logger, logErr *log.Logger, incoming chan tg.Message, outcoming chan tg.Message) {
	lastDialog, _ := os.OpenFile("lastDialog.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer lastDialog.Close()

	var wg sync.WaitGroup
	var consoleMutex sync.Mutex

	wg.Add(2)
	go ai.AiHeandler(login, password, logInfo, logErr, *lastDialog, &consoleMutex, incoming, outcoming, &wg)

	ctxTg, cancelTg := tg.SetupTg(logInfo, logErr, &consoleMutex, &wg)
	defer cancelTg()

	go tg.MonitorPartnerMessagesAndSend(ctxTg, logInfo, logErr, lastDialog, incoming, outcoming, 5*time.Second)

	wg.Wait()
	commandHandler(logInfo, logErr, incoming)

}

func main() {

	cfg := config.LoadConfig()
	os.Remove("QR.png")
	os.Remove("Chats.png")

	file, _ := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()

	os.Stderr = file

	logInfo := log.New(file, "APP_INFO\t", log.Ldate|log.Ltime)
	logErr := log.New(file, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	defer logInfo.Println("app exit")

	incoming := make(chan tg.Message, 10)
	outcoming := make(chan tg.Message, 1)

	dialog(cfg.Gmail, cfg.Password, logInfo, logErr, incoming, outcoming)
}
